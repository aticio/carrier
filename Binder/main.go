package main

import (
	"fmt"
	"net"

	server "carrier/Binder/domestic/binderserver"

	pb "carrier/Binder/rpc/binder"

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
	port := 8082
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("failed to listen")
	}

	grpcServer := grpc.NewServer()
	s, err := server.NewBinderServer()

	if err != nil {
		log.Errorf("cannot get new Binder grpc server: %s", err)
	}
	pb.RegisterServerServiceServer(grpcServer, s)
	log.Info("starting Binder grpc server on port ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to serve: %s", err)
	}
}
