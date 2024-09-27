package main

import (
	"context"
	"github.com/ndemeshchenko/zenith/pkg/components/api"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	"github.com/ndemeshchenko/zenith/pkg/components/heartbeat"
	l "github.com/ndemeshchenko/zenith/pkg/components/logger"
	"github.com/ndemeshchenko/zenith/pkg/components/mongo"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// set GOMAXPROCS
	_, _ = maxprocs.Set()
	appConfig := config.InitConfigurations()
	logLevel := appConfig.LogLevel
	l.Init(logLevel, nil)

	l.Logger.Info("initializing MongoDB connection")
	mongoClient, err := mongo.InitDBConnection(appConfig)
	if err != nil {
		panic(err)
	}

	// new Monitor instance
	monitor := heartbeat.NewMonitor(mongoClient)
	go monitor.Run()

	l.Logger.Info("starting Zenithd")
	api.Init(appConfig, mongoClient)

	defer func(mongoClient *mongo2.Client, ctx context.Context) {
		_ = mongoClient.Disconnect(ctx)
	}(mongoClient, context.Background())
}
