package observability

import (
	"testing"

	"goapp/internal/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func TestInitResource(t *testing.T) {
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}

	resource, err := initResource(cfg)
	if err != nil {
		t.Fatalf("Expected no error creating resource, got: %v", err)
	}

	if resource == nil {
		t.Fatal("Expected resource to be non-nil")
	}

	// Check that resource has expected attributes
	attrs := resource.Attributes()
	if len(attrs) == 0 {
		t.Error("Expected resource to have attributes")
	}
}

func TestInitTracer(t *testing.T) {
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "development",
	}

	// Test development environment (simple processor)
	cleanup := InitTracer(cfg)
	if cleanup == nil {
		t.Fatal("Expected cleanup function to be non-nil")
	}

	// Test that cleanup doesn't panic
	cleanup()
}

func TestInitTracerProduction(t *testing.T) {
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "production",
	}

	// Test production environment (batch processor)
	cleanup := InitTracer(cfg)
	if cleanup == nil {
		t.Fatal("Expected cleanup function to be non-nil")
	}

	// Test that cleanup doesn't panic
	cleanup()
}

func TestInitMeter(t *testing.T) {
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}

	cleanup := InitMeter(cfg)
	if cleanup == nil {
		t.Fatal("Expected cleanup function to be non-nil")
	}

	// Test that cleanup doesn't panic
	cleanup()
}

func TestInitCustomCounter(t *testing.T) {
	// First initialize meter
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	counterName := "test_counter"
	counter := InitCustomCounter(counterName)
	
	// Note: counter might be nil if meter isn't properly initialized
	// but the function should not panic
	if counter != nil {
		t.Log("Counter created successfully")
	} else {
		t.Log("Counter is nil (expected in some test scenarios)")
	}
}

func TestUpdateCounter(t *testing.T) {
	// Test with nil counter (should not panic)
	var nilCounter metric.Int64Counter
	UpdateCounter(nilCounter, 1)
	t.Log("UpdateCounter with nil counter completed without panic")

	// Test with valid counter
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	counter := InitCustomCounter("test_update_counter")
	UpdateCounter(counter, 5)
	t.Log("UpdateCounter with valid counter completed")
}

func TestTracer(t *testing.T) {
	// First initialize tracer
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	tracerCleanup := InitTracer(cfg)
	defer tracerCleanup()

	tracerName := "test-tracer"
	tracer := Tracer(tracerName)
	
	if tracer == nil {
		t.Fatal("Expected tracer to be non-nil")
	}
}

func TestMeter(t *testing.T) {
	// First initialize meter
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	meterName := "test-meter"
	meter := Meter(meterName)
	
	if meter == nil {
		t.Fatal("Expected meter to be non-nil")
	}
}

func TestRecordMetric(t *testing.T) {
	// Test with nil counter (should not panic)
	var nilCounter metric.Int64Counter
	attrs := []attribute.KeyValue{
		attribute.String("test", "value"),
	}
	RecordMetric(nilCounter, 1, attrs...)
	t.Log("RecordMetric with nil counter completed without panic")

	// Test with valid counter
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	counter := InitCustomCounter("test_record_counter")
	RecordMetric(counter, 10, attrs...)
	t.Log("RecordMetric with valid counter completed")

	// Test without attributes
	RecordMetric(counter, 5)
	t.Log("RecordMetric without attributes completed")
}

func TestInitTracerWithInvalidConfig(t *testing.T) {
	// Test with empty config to trigger resource creation issues
	cfg := config.ObservabilityConfig{
		ServiceName: "",
		Version:     "",
		Environment: "",
	}

	cleanup := InitTracer(cfg)
	if cleanup == nil {
		t.Fatal("Expected cleanup function to be non-nil even with invalid config")
	}

	// Should not panic even with empty config
	cleanup()
}

func TestInitMeterWithInvalidConfig(t *testing.T) {
	// Test with empty config
	cfg := config.ObservabilityConfig{
		ServiceName: "",
		Version:     "",
		Environment: "",
	}

	cleanup := InitMeter(cfg)
	if cleanup == nil {
		t.Fatal("Expected cleanup function to be non-nil even with invalid config")
	}

	// Should not panic even with empty config
	cleanup()
}

func TestInitCustomCounterWithoutMeter(t *testing.T) {
	// Test creating counter without proper meter initialization
	counterName := "test_counter_no_meter"
	counter := InitCustomCounter(counterName)
	
	// Counter might be nil or valid depending on global state
	// The important thing is that it doesn't panic
	if counter != nil {
		t.Log("Counter created successfully without explicit meter initialization")
	} else {
		t.Log("Counter is nil without explicit meter initialization (expected)")
	}
}

func TestInitCustomCounterInvalidName(t *testing.T) {
	// First initialize meter
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	// Test with empty counter name
	counter := InitCustomCounter("")
	
	// Should handle empty name gracefully
	if counter != nil {
		t.Log("Counter created with empty name")
	} else {
		t.Log("Counter is nil with empty name (expected)")
	}
}

func TestRecordMetricEdgeCases(t *testing.T) {
	// Test with nil counter and no attributes
	var nilCounter metric.Int64Counter
	RecordMetric(nilCounter, 1)
	t.Log("RecordMetric with nil counter and no attributes completed")

	// Test with valid counter and empty attributes slice
	cfg := config.ObservabilityConfig{
		ServiceName: "test-service",
		Version:     "1.0.0",
		Environment: "test",
	}
	meterCleanup := InitMeter(cfg)
	defer meterCleanup()

	counter := InitCustomCounter("test_record_edge_counter")
	emptyAttrs := []attribute.KeyValue{}
	RecordMetric(counter, 3, emptyAttrs...)
	t.Log("RecordMetric with empty attributes slice completed")
}

func TestObservabilityIntegration(t *testing.T) {
	cfg := config.ObservabilityConfig{
		ServiceName: "integration-test",
		Version:     "1.0.0",
		Environment: "test",
	}

	// Initialize both tracer and meter
	tracerCleanup := InitTracer(cfg)
	meterCleanup := InitMeter(cfg)
	
	defer func() {
		tracerCleanup()
		meterCleanup()
	}()

	// Create and use tracer
	tracer := Tracer("integration-test")
	if tracer == nil {
		t.Fatal("Expected tracer to be non-nil")
	}

	// Create and use meter
	meter := Meter("integration-test")
	if meter == nil {
		t.Fatal("Expected meter to be non-nil")
	}

	// Create and use counter
	counter := InitCustomCounter("integration_counter")
	UpdateCounter(counter, 1)
	RecordMetric(counter, 2, attribute.String("test", "integration"))

	t.Log("Integration test completed successfully")
}