package shared

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	mytrace "go.opentelemetry.io/otel/trace"
)

func InitTracing() {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}
	// Set sample rate to 1.0 for demonstration purposes.
	// In production, you might want to set a lower sample rate.

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// 1% is a common sample rate, but you can adjust it as needed.
		// trace.WithSampler(trace.TraceIDRatioBased(0.01)), // 1% sample rate
		trace.WithSampler(trace.AlwaysSample()), // Always sample for demo purposes
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		// Use from callers to propagate context
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func StartNewSpan(ctx context.Context, serviceName string, operationName string) (context.Context, mytrace.Span) {
	ctx, span := otel.Tracer(serviceName).Start(ctx, operationName)
	return ctx, span
}
