package main

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	// Initialize OpenTelemetry tracing
	cleanup, err := setupTraceProvider(os.Getenv("OTEL_ENDPOINT"), os.Getenv("SERVICE_NAME"), "1.0.0")
	if err != nil {
		log.Fatalf("Failed to set up trace provider: %v", err)
	}
	defer cleanup()

	// Initialize Prometheus metrics
	router := gin.New()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	// Register OpenTelemetry middleware
	router.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME")))

	// Initialize Slog logger
	// with span ID and trace ID enabled
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config := sloggin.Config{
		WithRequestID: true,
		WithSpanID:    true,
		WithTraceID:   true,
	}
	router.Use(sloggin.NewWithConfig(logger, config))

	router.Use(gin.Recovery())

	// Example route
	router.GET("/ping", func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		if rand.Intn(2) == 0 {
			sloggin.AddCustomAttributes(c, slog.String("message", "error"))
			c.JSON(500, gin.H{
				"message": "error",
			})
			return
		}
		sloggin.AddCustomAttributes(c, slog.String("ping", "pong"))
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/order", func(c *gin.Context) {
		var body struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		total := mockOrder(c, body)

		sloggin.AddCustomAttributes(c, slog.Any("total", total))
		c.JSON(200, gin.H{
			"message": "order",
			"total":   total,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}

func mockOrder(c *gin.Context, body struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}) float64 {
	sloggin.AddCustomAttributes(c, slog.Any("body", body))
	return float64(body.Quantity * 100)
}

func setupTraceProvider(endpoint string, serviceName string, serviceVersion string) (func(), error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	cleanup := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}
	return cleanup, nil
}
