package main

import (
	"context"
	"log"
	"os"
	"shared"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// createMessageSpans creates and configures tracing spans for message processing
func createMessageSpans(parentCtx context.Context, queueName string) (context.Context, []trace.Span) {
	spans := make([]trace.Span, 2)

	// Create consume span
	ctx, consumeSpan := otel.Tracer("service_b").Start(parentCtx, "consume-rabbitmq-message")
	consumeSpan.SetAttributes(
		attribute.String("queue.name", queueName),
	)
	spans[0] = consumeSpan

	// Create processing span
	_, processSpan := otel.Tracer("service_b").Start(ctx, "process-rabbitmq-message")
	spans[1] = processSpan

	return ctx, spans
}

// processMessageData handles the business logic of processing the message
func processMessageData(body []byte) {
	log.Printf(" [x] %s", body)
	log.Printf(" [x] Processing message: %s", body)
	time.Sleep(2 * time.Second)
	log.Printf(" [x] Finished processing message: %s", body)
}

func processData(headers map[string]string, body []byte, queueName string) {
	// Extract tracing context
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	parentCtx := propagator.Extract(context.Background(), propagation.MapCarrier(headers))

	// Create spans
	_, spans := createMessageSpans(parentCtx, queueName)
	// Ensure spans are closed
	for _, span := range spans {
		defer span.End()
	}

	log.Printf(" [x] Received a message with headers: %v", headers)

	// Execute business logic
	processMessageData(body)
}

func main() {
	shared.InitTracing()

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"users",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"users", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Extract headers
			headers := make(map[string]string)
			for k, v := range d.Headers {
				if str, ok := v.(string); ok {
					headers[k] = str
				}
			}
			processData(headers, d.Body, q.Name)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
