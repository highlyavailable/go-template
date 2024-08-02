package kafka

import "github.com/IBM/sarama"

// Config holds the configuration for Kafka connection
type Config struct {
	Brokers        []string `envconfig:"BROKERS" required:"true" split_words:"true"`
	ProducerTopic  string   `envconfig:"PRODUCER_TOPIC" required:"true"`
	ConsumerTopic  string   `envconfig:"CONSUMER_TOPIC" required:"true"`
	ConsumerGroup  string   `envconfig:"CONSUMER_GROUP" required:"true"`
	ConsumerOffset string   `envconfig:"CONSUMER_OFFSET" default:"oldest"`
}

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
