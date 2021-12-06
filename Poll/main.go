package main

import (
	pbscraper "carrier/Scraper/rpc/scraper"
	pbseed "carrier/Seed/rpc/seed"
	"context"
	"fmt"

	"errors"

	"github.com/go-test/deep"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	//seeds := getSeeds()
	initCompare()
	gocron.Every(15).Minutes().Do(initCompare)
	<-gocron.Start()
}

func getSeeds() []*pbseed.Seed {
	log.Info("getting all seeds from Seed service")

	port := 8081
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("cannot connect to Seed grpc server %v", err)
	}
	defer conn.Close()
	client := pbseed.NewSeedServiceClient(conn)

	result, err := client.GetAll(context.Background(), &pbseed.Void{})
	if err != nil {
		log.Errorf("error getting result form Seed grpc server %v", err)
	}
	return result.GetSeeds()
}

func initCompare() {
	seeds := getSeeds()
	for _, seed := range seeds {
		log.Info("comparing configuration for " + seed.GetHost() + ":" + seed.GetPort())
		triggerCompareRoutine(seed)
	}
}

func triggerCompareRoutine(seed *pbseed.Seed) {
	go compareSeeds(seed)
}

func compareSeeds(seed *pbseed.Seed) {
	scraperSeed, err := getConf(seed)
	if err != nil {
		log.Errorf("cannot check configuration updates for %s:%s", seed.Host, seed.Port)
		return
	}

	if diff := deep.Equal(scraperSeed, seed); diff == nil {
		log.Info(fmt.Sprintf("configuration is equal %s:%s", seed.GetHost(), seed.GetPort()))
	} else {
		log.Info(fmt.Sprintf("there are updates on liberty server %s:%s", seed.GetHost(), seed.GetPort()))
		updateSeed(*scraperSeed)
	}
}

func getConf(seed *pbseed.Seed) (*pbseed.Seed, error) {
	log.Info("getting configuration from Scraper service for " + seed.Host + ":" + seed.Port)

	port := 8083
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())

	if err != nil {
		log.Errorf("cannot connect to Seed grpc server %v", err)
	}
	defer conn.Close()
	client := pbscraper.NewScraperServiceClient(conn)

	scraperSeed, err := client.GetServerSeed(context.Background(), &pbscraper.Server{Host: seed.GetHost(), Port: seed.GetPort()})
	if err != nil {
		log.Errorf("Error getting result form Scraper grpc server %v", err)
		return nil, errors.New("Error getting result form Scraper grpc server ")
	}
	return scraperSeed, nil
}

func updateSeed(seed pbseed.Seed) {
	log.Info("updating seed " + seed.GetHost() + ":" + seed.GetPort())

	port := 8081
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())

	if err != nil {
		log.Errorf("cannot connect to Seed grpc server %v", err)
	}
	defer conn.Close()
	client := pbseed.NewSeedServiceClient(conn)

	result, err := client.UpdateSeed(context.Background(), &seed)
	if err != nil {
		log.Errorf("error getting result form Seed grpc server while update %v", err)
	}
	log.Info("update result of " + seed.GetHost() + ":" + seed.GetPort() + " is " + result.GetStatus())
}
