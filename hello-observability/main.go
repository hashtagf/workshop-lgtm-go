package main

import (
	"fmt"
	"net/http"
	"strconv"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	totalPingRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_get_ping_count",
			Help: "Total number of get api ping requests grouped by status",
		},
		[]string{"status"},
	)
	totalPingRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_get_ping_duration",
			Help: "Duration of get api ping requests",
		},
		[]string{"status"},
	)
)

func main() {
	fmt.Println("Hello, World!")
	startGinServer()
}

func startGinServer() {

	// custom metrics
	prometheus.MustRegister(totalPingRequest)
	prometheus.MustRegister(totalPingRequestDuration)

	// start gin server
	router := gin.Default()
	// Recovery gin
	router.Use(gin.Recovery())

	// add metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.GET("/ping", PingHandler)
	router.Run(":8080")
}

func PingHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	timer := prometheus.NewTimer(totalPingRequestDuration.WithLabelValues(strconv.Itoa(http.StatusOK)))
	if rand.Intn(2) == 0 {
		totalPingRequest.WithLabelValues(strconv.Itoa(http.StatusInternalServerError)).Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}
	timer.ObserveDuration()

	totalPingRequest.WithLabelValues(strconv.Itoa(http.StatusOK)).Inc()
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
