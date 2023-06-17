package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	alertModel "github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	"github.com/ndemeshchenko/zenith/pkg/components/models/environment"
	prometheusWebhook "github.com/ndemeshchenko/zenith/pkg/components/webhooks/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func Init(config *config.Config, mongoClient *mongo.Client) {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	//router.Use(cors.Default())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
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
				alerts, err := alertModel.GetAll(mongoClient)
				if err != nil {
					log.Printf(err.Error())
				}
				c.JSON(http.StatusOK, alerts)
			})
		}

		alert := v1.Group("/alert")
		{
			alert.GET("/:id", func(c *gin.Context) {
				alert, err := alertModel.GetOne(mongoClient, c.Param("id"))
				if err != nil {
					log.Printf(err.Error())
				}
				c.JSON(http.StatusOK, alert)
			})
			alert.PATCH("/:id", func(c *gin.Context) {
				//parse options
				action := c.Query("action")
				err := alertModel.UpdateStatus(mongoClient, c.Param("id"), action)
				if err != nil {
					log.Printf(err.Error())
				}
				c.JSON(http.StatusOK, gin.H{
					"status": "patched",
				})
			})
			alert.DELETE("/:id", func(c *gin.Context) {
				err := alertModel.DeleteOne(mongoClient, c.Param("id"))
				if err != nil {
					log.Printf(err.Error())
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
					fmt.Printf(err.Error())
				}
				c.JSON(http.StatusOK, environments)
			})
		}
	}

	router.Run("0.0.0.0:8080")
}

// TokenAuthMiddleware Super simple auth middleware just to cover the basics
func TokenAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		//log.Println("Authorization: Bearer ", authToken)
		if authToken != fmt.Sprintf("Bearer %s", token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		c.Next()
	}
}
