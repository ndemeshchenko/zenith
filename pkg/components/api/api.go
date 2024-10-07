package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	l "github.com/ndemeshchenko/zenith/pkg/components/logger"
	alertModel "github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	"github.com/ndemeshchenko/zenith/pkg/components/models/environment"
	heartbeatModel "github.com/ndemeshchenko/zenith/pkg/components/models/heartbeat"
	prometheusWebhook "github.com/ndemeshchenko/zenith/pkg/components/webhooks/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Init(config *config.Config, mongoClient *mongo.Client) {
	if config.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	//router.Use(cors.Default())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.FE_URL}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	// Register the middleware
	router.Use(cors.New(corsConfig))

	healthz := router.Group("/healthz")
	{
		healthz.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		})
	}

	v1 := router.Group("/api", TokenAuthMiddleware(config.AuthToken))
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/prometheus", func(c *gin.Context) {
				err := prometheusWebhook.ProcessWebhookAlert(c.Request.Body, mongoClient)
				if err != nil {
					l.Logger.Error(err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				c.JSON(200, gin.H{
					"alert-processed": true,
				})
			})
		}

		alerts := v1.Group("/alerts")
		{
			alerts.GET("", func(c *gin.Context) {
				//TODO add more complex query parsing
				filter := bson.M{
					"status": "firing",
					"type":   bson.M{"$ne": "heartbeat"},
				}
				envFilter := c.Query("environment")
				if envFilter != "" {
					filter["environment"] = envFilter
				}

				alerts, err := alertModel.GetAll(filter, mongoClient)
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, alerts)
			})
		}

		heartbeats := v1.Group("/heartbeats")
		{
			heartbeats.GET("", func(c *gin.Context) {
				filter := bson.M{}
				heartbeats, err := heartbeatModel.GetAll(filter, mongoClient)
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, heartbeats)
			})
		}

		alert := v1.Group("/alert")
		{
			alert.GET("/:id", func(c *gin.Context) {
				alert, err := alertModel.GetOne(mongoClient, c.Param("id"))
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, alert)
			})
			alert.PATCH("/:id", func(c *gin.Context) {
				//parse options
				action := c.Query("action")
				err := alertModel.UpdateStatus(mongoClient, c.Param("id"), action)
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, gin.H{
					"status": "patched",
				})
			})
			alert.DELETE("/:id", func(c *gin.Context) {
				err := alertModel.DeleteOne(mongoClient, c.Param("id"))
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, gin.H{
					"status": "deleted",
				})
			})
		}

		environments := v1.Group("/environments")
		{
			environments.GET("", func(c *gin.Context) {
				environments, err := environment.GetAll(mongoClient)
				if err != nil {
					l.Logger.Error(err.Error())
				}
				c.JSON(http.StatusOK, environments)
			})
		}
	}

	if config.EnableTLS && config.TLSCertFile != "" && config.TLSKeyFile != "" {
		l.Logger.Info("Starting Zenithd with TLS enabled")
		if err := router.RunTLS("0.0.0.0:443", config.TLSCertFile, config.TLSKeyFile); err != nil {
			panic(err)
		}
	} else {
		l.Logger.Info("Starting Zenithd :8080")
		if err := router.Run("0.0.0.0:8080"); err != nil {
			panic(err)
		}
	}
}

// TokenAuthMiddleware Super simple auth middleware just to cover the basics
func TokenAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken != fmt.Sprintf("Bearer %s", token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		c.Next()
	}
}
