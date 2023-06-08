package main

import (
	"context"
	"github.com/ndemeshchenko/zenith/pkg/components/api"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	"github.com/ndemeshchenko/zenith/pkg/components/mongo"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
)

func main() {
	log.Println("Starting Zenithd")
	// set GOMAXPROCS
	_, _ = maxprocs.Set()

	log.SetFlags(log.Ldate | log.Lshortfile)
	log.Println("Iitializing configurations")
	appConfig := config.InitConfigurations()

	log.Println("Iitializing MongoDB connection")
	mongoClient, error := mongo.InitDBConnection(appConfig)
	if error != nil {
		panic(error)
	}
	log.Println("authToken: ", appConfig.AuthToken)
	api.Init(appConfig, mongoClient)

	defer mongoClient.Disconnect(context.Background())
}

//X-Forwarded-Port: 80
//Host: zenith-monitor-1012986648.eu-central-1.elb.amazonaws.com
//X-Amzn-Trace-Id: Root=1-6481d4d7-4f7f4dc771cb34eb0f1719a6
//Content-Length: 2042
//User-Agent: Alertmanager/0.25.0
//Authorization: Bearer 52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649
//Content-Type: application/json
