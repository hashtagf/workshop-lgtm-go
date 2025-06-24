package main

import (
	"service_a/gateway"
	"service_a/user"
	"shared"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Initialize OpenTelemetry tracing
	shared.InitTracing()

	// Connect to RabbitMQ
	conn, err := shared.ConnectRabbitMQ()
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		panic("Failed to open a channel: " + err.Error())
	}
	defer channel.Close()
	err = channel.ExchangeDeclare(
		"users",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic("Failed to declare exchange: " + err.Error())
	}

	// Create a new Gin router
	r := gin.New()

	// Middleware for OpenTelemetry tracing
	// This will automatically instrument incoming HTTP requests
	r.Use(otelgin.Middleware("my-server"))

	r.Use(gin.Recovery())

	// Create a new user
	userHandler := user.UserHandler{Ch: channel}
	r.POST("/user", userHandler.CreateUser())

	// Get all products from HTTP request
	r.GET("/call/products", func(c *gin.Context) {
		ctx := c.Request.Context()
		products := gateway.CallGetAllProducts(ctx)
		c.JSON(200, gin.H{
			"message":  "Products fetched from service_c",
			"products": products,
		})
	})

	r.Run(":8080")
}
