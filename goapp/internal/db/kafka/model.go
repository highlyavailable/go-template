package kafka

import "github.com/IBM/sarama"

// Kafka holds the producer and consumer instances
type Kafka struct {
	Producer sarama.SyncProducer
	Consumer sarama.ConsumerGroup
}

// ConsumerGroupHandler is a custom implementation of the sarama.ConsumerGroupHandler interface
type ConsumerGroupHandler struct {
	handler func([]byte) error
	ready   chan bool
}
