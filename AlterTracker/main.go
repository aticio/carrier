package main

import (
	"fmt"

	"golang.org/x/net/context"

	dao "carrier/AlterTracker/domestic/dao"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	pbCarrier "carrier/Carrier/rpc/carrier"
	pbCarrierJvm "carrier/CarrierJvm/rpc/carrierjvm"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	readConfig()
	dao, err := dao.NewSeedDAO(viper.GetString("mongoUrl"))
	if err != nil {
		log.Errorf("could not ceate a database handler object, error %v", err)
	}

	log.Info("looking for changes...")
	coll := dao.GetColl()

	for {
		initWatching(coll)
	}
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

// func initCarrierClient() pbCarrier.CarrierServiceClient {
// 	port := 8084
// 	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

// 	if err != nil {
// 		log.Errorf("error connecting Seed RPC server for checking %v", err)
// 	}
// 	return pbCarrier.NewCarrierServiceClient(conn)
// }

func initWatching(coll *mgo.Collection) {
	pipeline := []bson.M{}

	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{FullDocument: mgo.UpdateLookup})
	if err == nil {
		var changeDoc map[string]interface{}
		for changeStream.Next(&changeDoc) {
			log.Info("In stream")
			if changeDoc["operationType"] == "replace" {
				replacePmi(changeDoc)
				replaceJvm(changeDoc)
			} else if changeDoc["operationType"] == "insert" {
				insertPmi(changeDoc)
				insertJvm(changeDoc)
			} else if changeDoc["operationType"] == "delete" {
				deletePmi(changeDoc)
				deleteJvm(changeDoc)
			}
		}

		if err := changeStream.Close(); err != nil {
			log.Info(err)
		}
	}
}

// PMI Changing
func replacePmi(changeDoc map[string]interface{}) {
	log.Info("replace")
	fd := changeDoc["fullDocument"].(map[string]interface{})
	port := 8084
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("error connecting Carrier RPC server for checking %v", err)
		return
	}
	client := pbCarrier.NewCarrierServiceClient(conn)
	result, err := client.UpdateCarrier(context.Background(), &pbCarrier.Carrier{Host: fd["host"].(string), Port: fd["port"].(string)})
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()
	log.Info(result.GetStatus())
}

func insertPmi(changeDoc map[string]interface{}) {
	log.Info("insert")
	fd := changeDoc["fullDocument"].(map[string]interface{})
	port := 8084
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))
	if err != nil {
		log.Errorf("error connecting Carrier RPC server for checking %v", err)
		return
	}
	client := pbCarrier.NewCarrierServiceClient(conn)
	result, err := client.AddCarrier(context.Background(), &pbCarrier.Carrier{Host: fd["host"].(string), Port: fd["port"].(string)})
	defer conn.Close()

	if err != nil {
		log.Error(err)
	}
	log.Info(result.GetStatus())
}

func deletePmi(changeDoc map[string]interface{}) {
	log.Info("delete")
	port := 8084
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))
	if err != nil {
		log.Errorf("error connecting Carrier RPC server for checking %v", err)
		return
	}
	client := pbCarrier.NewCarrierServiceClient(conn)
	result, err := client.DeleteCarrier(context.Background(), &pbCarrier.Void{})
	defer conn.Close()

	if err != nil {
		log.Error(err)
	}
	log.Info(result.GetStatus())
}

// JVM Changing

func replaceJvm(changeDoc map[string]interface{}) {
	log.Info("replace")
	fd := changeDoc["fullDocument"].(map[string]interface{})
	port := 8085
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))

	if err != nil {
		log.Errorf("error connecting CarrierJvm RPC server for checking %v", err)
		return
	}
	client := pbCarrierJvm.NewCarrierJvmServiceClient(conn)
	result, err := client.UpdateCarrierJvm(context.Background(), &pbCarrierJvm.CarrierJvm{Host: fd["host"].(string), Port: fd["port"].(string)})
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()
	log.Info(result.GetStatus())
}

func insertJvm(changeDoc map[string]interface{}) {
	log.Info("insert")
	fd := changeDoc["fullDocument"].(map[string]interface{})
	port := 8085
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))
	if err != nil {
		log.Errorf("error connecting CarrierJvm RPC server for checking %v", err)
		return
	}
	client := pbCarrierJvm.NewCarrierJvmServiceClient(conn)
	result, err := client.AddCarrierJvm(context.Background(), &pbCarrierJvm.CarrierJvm{Host: fd["host"].(string), Port: fd["port"].(string)})
	defer conn.Close()

	if err != nil {
		log.Error(err)
	}
	log.Info(result.GetStatus())
}

func deleteJvm(changeDoc map[string]interface{}) {
	log.Info("delete")
	port := 8085
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*16)))
	if err != nil {
		log.Errorf("error connecting CarrierJvm RPC server for checking %v", err)
		return
	}
	client := pbCarrierJvm.NewCarrierJvmServiceClient(conn)
	result, err := client.DeleteCarrierJvm(context.Background(), &pbCarrierJvm.Void{})
	defer conn.Close()

	if err != nil {
		log.Error(err)
	}
	log.Info(result.GetStatus())
}
