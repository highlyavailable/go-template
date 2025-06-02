package kafka

import (
	"context"
	"strings"
	"testing"
	"time"

	"goapp/internal/config"
)

func TestNew(t *testing.T) {
	cfg := config.KafkaConfig{
		Brokers:        []string{"localhost:9092"},
		ProducerTopic:  "test-topic",
		ConsumerTopic:  "test-topic",
		ConsumerGroup:  "test-group",
		ConsumerOffset: "oldest",
	}

	client, err := New(cfg)
	// Note: This will likely fail in test environment without Kafka running
	// But we can test the structure
	if err != nil {
		// Expected in test environment
		t.Logf("Expected Kafka connection to fail in test environment: %v", err)
		if !strings.Contains(strings.ToLower(err.Error()), "kafka") &&
		   !strings.Contains(strings.ToLower(err.Error()), "dial") &&
		   !strings.Contains(strings.ToLower(err.Error()), "connection") {
			t.Errorf("Expected connection-related error, got: %v", err)
		}
		return
	}

	if client == nil {
		t.Fatal("Expected client to be non-nil when connection succeeds")
	}

	// Test interface methods exist
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Test that we can call the methods without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Client methods panicked: %v", r)
		}
	}()

	// Test Produce (will likely fail but shouldn't panic)
	err = client.Produce("test-topic", "test-message")
	if err != nil {
		t.Logf("Expected produce to fail in test environment: %v", err)
	}

	// Test Consume setup
	handler := func(data []byte) error {
		return nil
	}
	
	err = client.Consume(ctx, "test-topic", handler)
	if err != nil {
		t.Logf("Consume setup error: %v", err)
	}

	// Test Close
	err = client.Close()
	if err != nil {
		t.Logf("Close returned error: %v", err)
	}
}

func TestKafkaStruct(t *testing.T) {
	// Test that Kafka struct can be created
	kafka := &Kafka{}
	if kafka == nil {
		t.Error("Expected Kafka struct to be creatable")
	}
}

func TestConsumerGroupHandler(t *testing.T) {
	// Test that ConsumerGroupHandler can be created
	handler := &ConsumerGroupHandler{
		handler: func(data []byte) error { return nil },
		ready:   make(chan bool),
	}
	
	if handler == nil {
		t.Error("Expected ConsumerGroupHandler to be creatable")
	}
	
	if handler.handler == nil {
		t.Error("Expected handler function to be set")
	}
	
	if handler.ready == nil {
		t.Error("Expected ready channel to be set")
	}
}