package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	// Custom metric counter_user_request by status code
	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_get_user_id_count",
			Help: "Number of requests of GET /user/:id endpoint by status",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(totalRequests)
	p := ginprometheus.NewPrometheus("gin")

	// Gin server
	router := gin.New()
	p.Use(router)
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		switch id {
		case "123":
			totalRequests.WithLabelValues("success").Inc()
			c.JSON(200, gin.H{
				"id":   id,
				"name": "User " + id,
			})
		case "456":
			totalRequests.WithLabelValues("error").Inc()
			c.JSON(200, gin.H{
				"error": "Internal server error",
			})
		default:
			totalRequests.WithLabelValues("not_found").Inc()
			c.JSON(200, gin.H{
				"error": "User not found",
			})
		}
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
