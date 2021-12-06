package explorerserver

import (
	pbBinder "carrier/Binder/rpc/binder"
	"carrier/Explorer/domestic/model"
	pbSeed "carrier/Seed/rpc/seed"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func checkSeed(server model.Server) bool {
	port := 8081
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("error connecting Seed RPC server for checking %v", err)
	}
	defer conn.Close()

	client := pbSeed.NewSeedServiceClient(conn)
	result, err := client.GetSeed(context.Background(), &pbSeed.Seed{Host: server.Host, Port: server.Port})
	if err != nil {
		log.Errorf("Seed RPC server didn't return any results for checking %v", err)
		return false
	}
	log.Info(fmt.Sprintf("Server %s:%s already added", result.GetHost(), result.GetPort()))
	return true
}

func sendServer(server model.Server) string {
	port := 8082
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())

	if err != nil {
		log.Errorf("error connecting Binder RPC server %v", err)
	}
	defer conn.Close()

	client := pbBinder.NewServerServiceClient(conn)
	result, err := client.HandleServer(context.Background(), &pbBinder.Server{Host: server.Host, Port: server.Port})
	if err != nil {
		log.Errorf("Binder RPC server didn't return any results %v", err)
		return result.GetStatus()
	}
	log.Info(result.GetStatus())
	return result.GetStatus()
}
