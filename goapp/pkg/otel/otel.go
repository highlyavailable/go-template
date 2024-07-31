package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/metric"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// initResource initializes the OpenTelemetry resource with common attributes.
func initResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("goapp"),
		attribute.String("service.version", "1.0.0"),
		attribute.String("environment", "development"),
	)
}

// InitMeter initializes the OpenTelemetry meter and returns a function that can be used to shutdown the meter.
func InitMeter() func() {
	exporter, err := prometheus.New(prometheus.WithoutUnits())
	if err != nil {
		log.Fatal(err)
	}

	meterProvider := sdk_metric.NewMeterProvider(
		sdk_metric.WithReader(exporter),
		sdk_metric.WithResource(initResource()),
	)
	otel.SetMeterProvider(meterProvider)

	return func() {
		_ = meterProvider.Shutdown(context.Background())
	}
}

// InitTracer initializes and configures the OpenTelemetry tracer for the application.
// It returns a function that can be used to shutdown the tracer provider.
func InitTracer() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(initResource()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		_ = tracerProvider.Shutdown(context.Background())
	}
}

// InitCustomCounter initializes a custom counter and returns it
func InitCustomCounter(counterName string) metric.Int64Counter {
	// Access the goapp meter
	meter := otel.GetMeterProvider().Meter("goapp")

	// Create a counter
	counter, err := meter.Int64Counter("custom_counter")
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}

	// Record a metric
	counter.Add(context.Background(), 0)

	return counter
}

func UpdateCounter(counter metric.Int64Counter, value int64) {
	// Record a metric
	counter.Add(context.Background(), value)
}
