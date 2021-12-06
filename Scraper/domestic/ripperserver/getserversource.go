package scraperserver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"

	pbSeed "carrier/Seed/rpc/seed"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//MetricConf struct
type MetricConf struct {
	Name     string   `yaml:"name"`
	Regex    string   `yaml:"regex"`
	URL      string   `yaml:"url"`
	Interval string   `yaml:"interval"`
	Metrics  []Metric `yaml:"metrics"`
}

//Metric struct
type Metric struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Attribute struct {
	MetricType string
	Name       string
	URL        string
}

func getMBeans(host, port string) (string, error) {
	log.Info("getting mbeans from ", host, port)
	req, err := http.NewRequest("GET", "https://"+host+":"+port+viper.GetString("libertyMBeansPath"), nil)
	fmt.Println("https://" + host + ":" + port + viper.GetString("libertyMBeansPath"))
	if err == nil {
		req.SetBasicAuth(viper.GetString("libertyUsername"), viper.GetString("libertyPassword"))
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: tr}
		resp, err := httpClient.Do(req)

		if err == nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			s := buf.String()
			if len(s) > 0 {
				if string(s[0:9]) != "<!DOCTYPE" {
					return s, nil
				}
				return "", errors.New("montitor feature has not been activated for " + host)
			}
			return "", errors.New("no response from server, check credentials " + host)
		}
		log.Error("error occured whilte getting response from server ", host, port)
		return "", err
	}
	log.Error("could not obtain http request for getting data ", host, port)
	return "", err
}

//ReadYAML reading yaml file which includes jsonpaths for ripping liberty server's json file
func readYAML() ([]MetricConf, error) {
	log.Info("reading yaml configuration...")
	yamlFile, err := ioutil.ReadFile(viper.GetString("metricConfPath"))
	if err != nil {
		log.Errorf("error reading yaml file %v ", err)
		return nil, err
	}

	metricConfs := []MetricConf{}

	err = yaml.Unmarshal(yamlFile, &metricConfs)
	if err != nil {
		log.Errorf("error unmarshalling yaml: %v", err)
		return nil, err
	}
	return metricConfs, nil
}

func RequestMetrics(metricConfs []MetricConf, mBeans, host, port string) *pbSeed.Seed {
	seed := pbSeed.Seed{}
	for _, metricConf := range metricConfs {
		attributes := getUrlAndName(metricConf, mBeans)
		for _, attribute := range attributes {
			pbAttribute := pbSeed.Seed_Attribute{Type: attribute.MetricType, Name: attribute.Name, Url: attribute.URL, Interval: metricConf.Interval}
			for _, metric := range metricConf.Metrics {
				m := pbSeed.Seed_Attribute_Metric{Name: metric.Name, Jpath: metric.Value}
				pbAttribute.Metrics = append(pbAttribute.Metrics, &m)
			}
			seed.Attributes = append(seed.Attributes, &pbAttribute)
		}
	}
	return &seed
}

func getUrlAndName(metricConf MetricConf, mBeans string) []Attribute {
	v := interface{}(nil)
	json.Unmarshal([]byte(mBeans), &v)
	urlInterface, err := jsonpath.Get(metricConf.URL, v)
	if err != nil {
		log.Errorf("error following jsonpath %v: ", err)
	}

	var attributes []Attribute
	for _, e := range urlInterface.([]interface{}) {
		objectName := e.(map[string]interface{})["objectName"]
		className := e.(map[string]interface{})["className"].(string)

		classNameSplit := strings.Split(className, ".")
		metricType := classNameSplit[len(classNameSplit)-1]
		name := assembleName(objectName, metricConf)
		if metricConf.Regex != "" {
			var validName = regexp.MustCompile(metricConf.Regex)
			rx := validName.FindAllString(name, -1)
			rxs := strings.Join(rx, "")
			if rxs != "" {
				name = strings.Replace(name, rxs, "", -1)
			}
		}
		url := e.(map[string]interface{})["URL"].(string) + viper.GetString("metricContextSuffix")
		attribute := Attribute{MetricType: metricType, Name: name, URL: url}
		attributes = append(attributes, attribute)
	}
	return attributes
}

func assembleName(objectName interface{}, metricConf MetricConf) string {
	couples := strings.Split(objectName.(string), ",")
	keys := strings.Split(metricConf.Name, ",")

	name := ""
	for _, couple := range couples {
		l := strings.Split(couple, "=")
		for _, key := range keys {
			if key == l[0] {
				name += strings.Replace(l[1], "\"", "", -1) + "_"
			}
		}
	}
	name = strings.Replace(name, " ", "_", -1)
	return strings.TrimRight(name, "_")
}

func getJvmName(mBeans, host, port string) (string, error) {
	jvmJsonPath := viper.GetString("jvmJsonPath")
	jvmJsonPath2 := viper.GetString("jvmJsonPath2")
	v := interface{}(nil)
	json.Unmarshal([]byte(mBeans), &v)
	urlInterface, err := jsonpath.Get(jvmJsonPath, v)
	if err != nil {
		log.Errorf("error following jsonpath %v: ", err)
		return "", err
	}
	if len(urlInterface.([]interface{})) == 0 {
		urlInterface, err = jsonpath.Get(jvmJsonPath2, v)
		log.Info("Using application name for jvm: ", host, ":", port)
		if err != nil {
			log.Errorf("error following jsonpath %v: ", err)
			return "", err
		}
	}
	objectName := urlInterface.([]interface{})[0]
	p := strings.Split(objectName.(string), ",")
	c := strings.Split(p[1], "=")
	return c[1], nil
}
