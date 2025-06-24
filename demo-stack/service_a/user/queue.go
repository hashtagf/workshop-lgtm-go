package user

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func PublishUserCreationEvent(user User, ctx context.Context, ch *amqp091.Channel) {
	// Create spans and get updated context
	ctx, spans := createUserEventSpans(ctx, user)
	// Ensure spans are closed
	for _, span := range spans {
		defer span.End()
	}

	// Execute business logic
	if err := publishUserEvent(ctx, user, ch); err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Println("User creation event published successfully")
	}
}

func createUserEventSpans(ctx context.Context, user User) (context.Context, []trace.Span) {
	spans := make([]trace.Span, 2)

	// Serialize span
	spans[0] = startNewSpanWithAttributes(ctx, "SerializedUserCreationEvent",
		attribute.Int("user_id", user.ID),
		attribute.String("user_name", user.Name),
		attribute.String("user_email", user.Email),
		attribute.String("operation", "serialize"),
	)

	// Publish span
	spans[1] = startNewSpanWithAttributes(ctx, "PublishUserCreationEvent",
		attribute.String("exchange_name", "users"),
		attribute.String("operation", "publish"),
	)

	return ctx, spans
}

// publishUserEvent is a pure function that handles the business logic of publishing a user event
func publishUserEvent(ctx context.Context, user User, ch *amqp091.Channel) error {
	// Prepare headers for tracing
	headers := make(map[string]string)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator.Inject(ctx, propagation.MapCarrier(headers))

	// Marshal user data
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Convert headers to AMQP table
	amqpHeaders := make(amqp091.Table)
	for k, v := range headers {
		amqpHeaders[k] = v
	}

	// Publish the message
	return ch.PublishWithContext(ctx,
		"users", // exchange
		"",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        userJSON,
			Headers:     amqpHeaders,
		})
}

func startNewSpan(ctx context.Context, operationName string) trace.Span {
	_, span := otel.Tracer("service-a").Start(ctx, operationName)
	return span
}

func startNewSpanWithAttributes(ctx context.Context, operationName string, attributes ...attribute.KeyValue) trace.Span {
	span := startNewSpan(ctx, operationName)
	span.SetAttributes(attributes...)
	return span
}
