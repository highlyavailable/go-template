package kafka

import (
	"context"
	"log"

	"goapp/internal/config"
)

// ExampleNew demonstrates how to use the New function to create a Kafka client,
// connect to a Kafka broker, produce a test message, and consume messages from a specified topic.
func ExampleNew() {
	cfg := config.KafkaConfig{
		Brokers:        []string{"localhost:9092"},
		ProducerTopic:  "events",
		ConsumerTopic:  "events", 
		ConsumerGroup:  "goapp-group",
		ConsumerOffset: "oldest",
	}

	kafkaClient, err := New(cfg)
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
