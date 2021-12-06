package binderserver

import (
	pbBinder "carrier/Binder/rpc/binder"
	pbScraper "carrier/Scraper/rpc/scraper"
	pbSeed "carrier/Seed/rpc/seed"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func getServerSeed(r *pbBinder.Server) (*pbSeed.Seed, error) {
	log.Info("getting server information for: ", r.GetHost(), r.GetPort())
	port := 8083
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("error connecting Scraper RPC server %v", err)
	}
	defer conn.Close()

	client := pbScraper.NewScraperServiceClient(conn)
	result, err := client.GetServerSeed(context.Background(), &pbScraper.Server{Host: r.GetHost(), Port: r.GetPort()})
	if err != nil {
		log.Errorf("Scraper RPC server didn't return any results %v", err)
		return nil, err
	}
	seedAttributes := []*pbSeed.Seed_Attribute{}
	for _, attribute := range result.GetAttributes() {
		seedAttribute := &pbSeed.Seed_Attribute{
			Type:     attribute.GetType(),
			Name:     attribute.GetName(),
			Url:      attribute.GetUrl(),
			Interval: attribute.GetInterval(),
		}
		seedMetrics := []*pbSeed.Seed_Attribute_Metric{}
		for _, metric := range attribute.GetMetrics() {
			seedMetric := &pbSeed.Seed_Attribute_Metric{
				Name:  metric.GetName(),
				Jpath: metric.GetJpath(),
			}
			seedMetrics = append(seedMetrics, seedMetric)
		}
		seedAttribute.Metrics = seedMetrics
		seedAttributes = append(seedAttributes, seedAttribute)
	}

	return &pbSeed.Seed{
		Host:       result.GetHost(),
		Port:       result.GetPort(),
		Jvm:        result.GetJvm(),
		Username:   result.GetUsername(),
		Password:   result.GetPassword(),
		Attributes: seedAttributes,
	}, nil
}

func addSeed(seed *pbSeed.Seed) bool {
	log.Info("sending seed for adding to db: ", seed.GetHost(), seed.GetPort())
	port := 8081
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("error connecting Seed RPC server %v", err)
		return false
	}
	defer conn.Close()
	client := pbSeed.NewSeedServiceClient(conn)
	_, err = client.AddSeed(context.Background(), seed)
	if err != nil {
		log.Errorf("Seed RPC server didn't return any results %v", err)
		return false
	}
	return true
}
