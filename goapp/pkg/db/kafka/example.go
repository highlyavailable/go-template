package kafka

import (
	"context"
	"log"
)

// ExampleNewConfig demonstrates how to use the NewConfig function to create a Kafka configuration,
// connect to a Kafka broker, produce a test message, and consume messages from a specified topic.
func ExampleNewConfig() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("Could not load Kafka config: %v", err)
	}

	kafkaClient, err := Connect(cfg)
	if err != nil {
		log.Fatalf("Could not connect to Kafka: %v", err)
	}
	defer kafkaClient.Close()

	err = kafkaClient.Produce(cfg.ProducerTopic, "Test message")
	if err != nil {
		log.Fatalf("Could not produce message: %v", err)
	}

	log.Println("Message sent successfully")

	// Define the message handler
	handler := func(msg []byte) error {
		log.Printf("Received message: %s", string(msg))
		return nil
	}

	// Start consuming messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = kafkaClient.Consume(ctx, cfg.ConsumerTopic, handler)
	if err != nil {
		log.Fatalf("Could not consume messages: %v", err)
	}
}
