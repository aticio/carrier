package main

import (
	"fmt"
	"net"

	b "carrier/CarrierJvm/domestic/business"
	server "carrier/CarrierJvm/domestic/carrierserver"
	pb "carrier/CarrierJvm/rpc/carrierjvm"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	readConfig()
	go initService()
	b.InitOperations()
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func initService() {
	port := 8085
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("failed to listen")
	}

	grpcServer := grpc.NewServer()
	s, err := server.NewCarrierJvmServer()

	if err != nil {
		log.Errorf("cannot get new Carrier grpc server: %s", err)
	}

	pb.RegisterCarrierJvmServiceServer(grpcServer, s)
	log.Info("starting Carrier grpc server on port ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to serve: %s", err)
	}
}
