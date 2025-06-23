package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	sloggin "github.com/samber/slog-gin"
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
	router := gin.New()

	// add metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config := sloggin.Config{
		// WithRequestID: true,
		WithSpanID:  true,
		WithTraceID: true,
	}
	router.Use(sloggin.NewWithConfig(logger, config))

	router.Use(gin.Recovery())

	router.GET("/ping", PingHandler)
	router.Run(":8080")
}

func PingHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	timer := prometheus.NewTimer(totalPingRequestDuration.WithLabelValues(strconv.Itoa(http.StatusOK)))
	if rand.Intn(2) == 0 {
		totalPingRequest.WithLabelValues(strconv.Itoa(http.StatusInternalServerError)).Inc()
		sloggin.AddCustomAttributes(c, slog.String("message", "error"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}
	timer.ObserveDuration()

	totalPingRequest.WithLabelValues(strconv.Itoa(http.StatusOK)).Inc()

	sloggin.AddCustomAttributes(c, slog.String("message", "pong"))
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
