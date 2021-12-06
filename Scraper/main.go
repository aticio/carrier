package main

import (
	"fmt"
	"net"

	server "carrier/Scraper/domestic/scraperserver"

	pb "carrier/Scraper/rpc/scraper"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	readConfig()
	initService()
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func initService() {
	port := 8083
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s, err := server.NewScraperServer()
	if err != nil {
		log.Errorf("cannot get new Scraper grpc server: %s", err)
	}
	pb.RegisterScraperServiceServer(grpcServer, s)

	log.Info("starting Scraper grpc server on port ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to server: %s", err)
	}
}
