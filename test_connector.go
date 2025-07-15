package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	ctx := context.Background()

	// Configure OTLP HTTP exporter (send to local Collector)
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Failed to create OTLP trace exporter: ", err)
	}

	// Configure resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatal("Failed to create resource: ", err)
	}

	// Configure trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal("Failed to shutdown TracerProvider: ", err)
		}
	}()

	otel.SetTracerProvider(tp)

	// Create tracer
	tracer := otel.Tracer("test-tracer")

	// Test case 1: span with "request.n" attribute (connector should react)
	log.Println("🧪 Test 1: Sending trace with 'request.n' attribute...")
	_, span1 := tracer.Start(ctx, "test-operation-with-trigger")
	span1.SetAttributes(
		attribute.String("request.n", "some-value"), // This should trigger the connector
		attribute.String("http.method", "GET"),
		attribute.String("http.url", "/api/test"),
	)
	time.Sleep(5 * time.Second)
	span1.End()

	// Test case 2: span without "request.n" attribute (connector should not react)
	log.Println("🧪 Test 2: Sending trace without 'request.n' attribute...")
	_, span2 := tracer.Start(ctx, "test-operation-without-trigger")
	span2.SetAttributes(
		attribute.String("http.method", "POST"),
		attribute.String("http.url", "/api/other"),
		attribute.String("user.id", "12345"),
	)
	time.Sleep(5 * time.Second)
	span2.End()

	// Wait for all spans to be sent
	time.Sleep(2 * time.Second)
	log.Println("✅ Test completed! Please check the Collector logs.")
}
