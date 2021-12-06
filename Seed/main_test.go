package main

import (
	pb "carrier/Seed/rpc/seed"
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
)

func TestReadConfig(t *testing.T) {
	readConfig()
}

func TestInitService(t *testing.T) {
	go initService()
	port := 8080
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())

	if err != nil {
		t.Errorf("Cannot connect to grpc server %v", err)
	}

	client := pb.NewSeedServiceClient(conn)
	result, err := client.GetSeed(context.Background(), &pb.Seed{Host: "my_host", Port: "9443"})
	if err != nil {
		t.Errorf("Error getting result from grpc server %v", err)
	}
	t.Log(result)
	fmt.Println(result)
}
