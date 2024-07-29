package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
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

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(initResource()),
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
