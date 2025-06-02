package observability

import (
	"context"
	"fmt"

	"goapp/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	api_trace "go.opentelemetry.io/otel/trace"
)

// initResource creates an OpenTelemetry resource with service information
func initResource(cfg config.ObservabilityConfig) (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.Version),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
}

// InitTracer initializes OpenTelemetry tracing with proper error handling
func InitTracer(cfg config.ObservabilityConfig) func() {
	ctx := context.Background()
	
	res, err := initResource(cfg)
	if err != nil {
		// In enterprise systems, we don't panic - we return the error or use a fallback
		fmt.Printf("Failed to create resource: %v\n", err)
		return func() {} // Return no-op function
	}

	var exporter trace.SpanExporter
	
	// Use stdout exporter for now (can be extended to support OTLP later)
	stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		fmt.Printf("Failed to create stdout exporter: %v\n", err)
		return func() {}
	}
	exporter = stdoutExporter

	// Configure span processor based on environment
	var processor trace.SpanProcessor
	if cfg.Environment == "production" {
		// Batch processor for production (better performance)
		processor = trace.NewBatchSpanProcessor(exporter)
	} else {
		// Simple processor for development (immediate export for debugging)
		processor = trace.NewSimpleSpanProcessor(exporter)
	}

	// Create tracer provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(processor),
		trace.WithSampler(trace.AlwaysSample()), // Configure sampling as needed
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)
	
	// Set global text map propagator for distributed tracing
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}
}

// InitMeter initializes OpenTelemetry metrics with proper error handling
func InitMeter(cfg config.ObservabilityConfig) func() {
	ctx := context.Background()
	
	res, err := initResource(cfg)
	if err != nil {
		fmt.Printf("Failed to create resource: %v\n", err)
		return func() {}
	}

	// Create Prometheus exporter
	exporter, err := prometheus.New(
		prometheus.WithoutUnits(),
		prometheus.WithoutCounterSuffixes(),
	)
	if err != nil {
		fmt.Printf("Failed to create prometheus exporter: %v\n", err)
		return func() {}
	}

	// Create meter provider
	meterProvider := sdk_metric.NewMeterProvider(
		sdk_metric.WithReader(exporter),
		sdk_metric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	return func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down meter provider: %v\n", err)
		}
	}
}

// InitCustomCounter creates a named counter with proper error handling
func InitCustomCounter(counterName string) metric.Int64Counter {
	meter := otel.GetMeterProvider().Meter("goapp")
	
	counter, err := meter.Int64Counter(
		counterName,
		metric.WithDescription(fmt.Sprintf("Counter for %s", counterName)),
	)
	if err != nil {
		fmt.Printf("Failed to create counter %s: %v\n", counterName, err)
		// Return nil and handle in UpdateCounter
		return nil
	}
	
	return counter
}

// UpdateCounter increments a counter with proper error handling
func UpdateCounter(counter metric.Int64Counter, value int64) {
	if counter == nil {
		return
	}
	counter.Add(context.Background(), value)
}


// Tracer returns a tracer for the given name
func Tracer(name string) api_trace.Tracer {
	return otel.Tracer(name)
}

// Meter returns a meter for the given name
func Meter(name string) metric.Meter {
	return otel.Meter(name)
}

// RecordMetric is a helper function to record metrics with error handling
func RecordMetric(counter metric.Int64Counter, value int64, attrs ...attribute.KeyValue) {
	if counter == nil {
		return
	}
	
	options := make([]metric.AddOption, 0, len(attrs))
	if len(attrs) > 0 {
		options = append(options, metric.WithAttributes(attrs...))
	}
	
	counter.Add(context.Background(), value, options...)
}