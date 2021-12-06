package business

import (
	pbseed "carrier/Seed/rpc/seed"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/fsnotify/fsnotify"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/patrickmn/go-cache"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	c        *cache.Cache
	h        *http.Client
	keys     []string
	w        *kafka.Writer
	m        *gocron.Scheduler
	secs     uint64
	jvmCheck string
)

func InitOperations() {
	seedClient := initSeedClient()
	h = initHTTPClient()
	c = cache.New(cache.DefaultExpiration, 0)

	secs = uint64(viper.GetInt("tickInSeconds"))
	kafkaVMsString := viper.GetString("kafkaVMs")
	kafkaVMs := strings.Split(kafkaVMsString, ",")
	seeds := getSeedConfigs(seedClient)
	jvmCheck = viper.GetString("jvmCheckContext")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		secs = uint64(viper.GetInt("tickInSeconds"))
		kafkaVMsString = viper.GetString("kafkaVMs")
		jvmCheck = viper.GetString("jvmCheckContext")
		initIndexing()
	})

	for _, seed := range seeds {
		c.SetDefault(seed.GetHost()+":"+seed.GetPort(), seed)
	}

	keys = make([]string, 0, len(c.Items()))
	for k := range c.Items() {
		keys = append(keys, k)
	}
	w = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaVMs,
		Topic:    "custommon",
		Balancer: &kafka.LeastBytes{},
	})
	initIndexing()
}

func initSeedClient() pbseed.SeedServiceClient {
	log.Info("initiating Seed service grpc client...")
	port := 8081
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("Cannot connect to grpc server %v", err)
	}

	client := pbseed.NewSeedServiceClient(conn)
	return client
}

func getSeedConfigs(client pbseed.SeedServiceClient) []*pbseed.Seed {
	log.Info("getting seed configs from Seed service...")
	pbSeeds, err := client.GetAll(context.Background(), &pbseed.Void{})
	if err != nil {
		log.Errorf("Error getting result form Seed grpc server %v", err)
	}
	return pbSeeds.GetSeeds()
}

func initHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        250,
		MaxIdleConnsPerHost: 250,
		IdleConnTimeout:     35,
	}
	httpClient := &http.Client{Transport: tr, Timeout: time.Duration(35 * time.Second)}
	return httpClient
}

//Indexing
func initIndexing() {
	gocron.Clear()
	gocron.Every(secs).Seconds().Do(start)
	log.Info("new cron created")
	<-gocron.Start()
}

func start() {
	for _, key := range keys {
		if x, found := c.Get(key); found {
			seed := x.(*pbseed.Seed)
			for _, attribute := range seed.GetAttributes() {
				go initLeave(seed.GetHost(), seed.GetPort(), seed.GetJvm(), seed.GetUsername(), seed.GetPassword(), attribute)
			}
		}
	}
}

func initLeave(host, port, jvm, username, password string, attribute *pbseed.Seed_Attribute) {
	content, err := requestMetrics(host, port, attribute.GetUrl(), username, password)
	if err == nil {
		metrics := getMetrics(content, attribute)

		tags := map[string]string{"host": host, "port": port, "jvm": jvm, "name": attribute.GetName()}
		fields := map[string]interface{}{}

		for k, v := range metrics {
			if strings.Contains(v, ".") {
				m, err := strconv.ParseFloat(v, 64)
				if err == nil {
					fields[k] = m
				}
			} else {
				m, err := strconv.ParseInt(v, 10, 64)
				if err == nil {
					fields[k] = m
				}
			}
		}

		pt, err := client.NewPoint("liberty_"+attribute.GetType(), tags, fields, time.Now())

		if err == nil {
			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Value: []byte(pt.String()),
				},
			)
		} else {
			log.Error(err)
		}
	} else {
		log.Error(err)
	}
}

func requestMetrics(host, port, context, username, password string) (string, error) {
	url := "https://" + host + ":" + port + context
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Error("could not obtain http request for getting data: " + url)
		return "", err
	}
	req.Close = true
	req.SetBasicAuth(username, password)

	resp, err := h.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Error("error occured whilte getting response: error after h.Do(req) " + url)
		return "", err
	}

	if resp.StatusCode != 200 {
		log.Error("error occured while getting response: status error " + host + ":" + port + " " + resp.Status + "url: " + url)
		return "", errors.New("error occured whilte getting response: status error " + host + ":" + port + " " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error occured while getting response: io error " + host + ":" + port + " " + resp.Status + "url: " + url)
		return "", err
	}

	s := string(body)
	return s, nil
}

func getMetrics(content string, attribute *pbseed.Seed_Attribute) map[string]string {
	metrics := make(map[string]string)
	v := interface{}(nil)
	err := json.Unmarshal([]byte(content), &v)
	if err != nil {
		log.Errorf("error unmarshaling json: %s", err)
	}

	for _, metric := range attribute.GetMetrics() {
		metricInterface, err := jsonpath.Get(metric.GetJpath(), v)
		if err != nil {
			log.Errorf("error following jsonpath: %s", err)
		}
		ms := fmt.Sprintf("%v", metricInterface)
		m := ms[1 : len(ms)-1]
		metrics[metric.GetName()] = m
	}
	return metrics
}

func AddSeed(host, port string) error {
	seedClient := initSeedClient()
	result, err := seedClient.GetSeed(context.Background(), &pbseed.Seed{Host: host, Port: port})
	if err != nil {
		log.Error(err)
		return err
	}
	c.SetDefault(result.GetHost()+":"+result.GetPort(), result)
	keys = append(keys, result.GetHost()+":"+result.GetPort())
	return nil
}

func UpdateSeed(host, port string) error {
	seedClient := initSeedClient()
	result, err := seedClient.GetSeed(context.Background(), &pbseed.Seed{Host: host, Port: port})
	if err != nil {
		log.Error(err)
		return err
	}
	c.SetDefault(result.GetHost()+":"+result.GetPort(), result)
	return nil
}

func DeleteSeed() error {
	seedClient := initSeedClient()
	result, err := seedClient.GetAll(context.Background(), &pbseed.Void{})
	if err != nil {
		log.Error(err)
		return err
	}
	c.Flush()
	for _, seed := range result.GetSeeds() {
		c.SetDefault(seed.GetHost()+":"+seed.GetPort(), seed)
	}

	keys = keys[:0]
	for k := range c.Items() {
		keys = append(keys, k)
	}
	return nil
}
