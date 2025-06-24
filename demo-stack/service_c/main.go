package main

import (
	"log"
	"service_c/product"
	"shared"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Initialize OpenTelemetry tracing
	shared.InitTracing()

	// Connect to MongoDB
	client, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Gin router
	r := gin.New()

	// Middleware for OpenTelemetry tracing
	// This will automatically instrument incoming HTTP requests
	r.Use(otelgin.Middleware("gin-server"))

	r.Use(gin.Recovery())

	// Get /products endpoint
	h := &product.ProductHandler{
		Client: client,
	}
	r.GET("/products", h.GetAllProducts())

	r.Run(":8080")
}
