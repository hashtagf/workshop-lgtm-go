package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"service_c/product"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func CallGetAllProducts(ctx context.Context) []product.Product {
	return httpGet(ctx, "/products")
}

func httpGet(ctx context.Context, uri string) []product.Product {
	url := os.Getenv("PRODUCTS_URL") + uri

	// Create a new span for the HTTP request
	ctx, span := otel.Tracer("service_a").Start(ctx, "httpGet")
	defer span.End()
	span.SetAttributes(attribute.Key("url").String(url))

	// Use otelhttp to instrument the HTTP request
	resp, err := otelhttp.Get(ctx, url)
	if err != nil {
		fmt.Println("Failed to make HTTP request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Decode the response body into a slice of products
	var products []product.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		fmt.Println("Failed to decode response:", err)
		return nil
	}
	return products
}
