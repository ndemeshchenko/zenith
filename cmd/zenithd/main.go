package main

import (
	"context"
	"github.com/ndemeshchenko/zenith/pkg/components/api"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	"github.com/ndemeshchenko/zenith/pkg/components/heartbeat"
	"github.com/ndemeshchenko/zenith/pkg/components/mongo"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
)

func main() {
	log.Println("Starting Zenithd")
	// set GOMAXPROCS
	_, _ = maxprocs.Set()

	log.SetFlags(log.Ldate | log.Lshortfile)
	log.Println("Initializing configurations")
	appConfig := config.InitConfigurations()

	log.Println("Initializing MongoDB connection")
	mongoClient, err := mongo.InitDBConnection(appConfig)
	if err != nil {
		panic(err)
	}

	// new Monitor instance
	monitor := heartbeat.NewMonitor(mongoClient)
	go monitor.Run()

	//log.Println("authToken: ", appConfig.AuthToken)
	api.Init(appConfig, mongoClient)

	defer mongoClient.Disconnect(context.Background())
}
