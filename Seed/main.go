package main

import (
	"fmt"
	"net"

	server "carrier/Seed/domestic/seedserver"
	pb "carrier/Seed/rpc/seed"

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
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("failed to listen")
	}

	grpcServer := grpc.NewServer()
	s, err := server.NewSeedServer()

	if err != nil {
		log.Errorf("cannot get new Seed grpc server: %s", err)
	}

	pb.RegisterSeedServiceServer(grpcServer, s)
	log.Info("starting Seed grpc server on port ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to serve: %s", err)
	}
}
