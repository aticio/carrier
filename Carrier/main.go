package main

import (
	"fmt"
	"net"

	b "carrier/Carrier/domestic/business"
	server "carrier/Carrier/domestic/carrierserver"
	pb "carrier/Carrier/rpc/carrier"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	readConfig()
	// tick := uint64(viper.GetInt("tickInSeconds"))
	// kafkaVMsString := viper.GetString("kafkaVMs")
	// jvmCheckContext := viper.GetString("jvmCheckContext")
	go initService()
	//b.InitOperations(tick, kafkaVMsString, jvmCheckContext)
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
	port := 8084
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("failed to listen")
	}

	grpcServer := grpc.NewServer()
	s, err := server.NewCarrierServer()

	if err != nil {
		log.Errorf("cannot get new Carrier grpc server: %s", err)
	}

	pb.RegisterCarrierServiceServer(grpcServer, s)
	log.Info("starting Carrier grpc server on port ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("failed to serve: %s", err)
	}
}
