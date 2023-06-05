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
	// set GOMAXPROCS
	_, _ = maxprocs.Set()

	log.SetFlags(log.Ldate | log.Lshortfile)

	appConfig := config.InitConfigurations()
	mongoClient, error := mongo.InitDBConnection(appConfig)
	if error != nil {
		panic(error)
	}
	api.Init(appConfig, mongoClient)

	defer mongoClient.Disconnect(context.Background())
}
