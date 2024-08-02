package kafka

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
)

// NewConfig creates a new Kafka configuration from environment variables
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("KAFKA", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}
	return &cfg, nil
}

// Connect initializes the Kafka producer and consumer
func Connect(cfg *Config) (*Kafka, error) {
	kafka := &Kafka{}

	// Initialize Producer
	producer, err := newProducer(cfg.Brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	kafka.Producer = producer

	// Initialize Consumer
	consumer, err := newConsumer(cfg.Brokers, cfg.ConsumerGroup, cfg.ConsumerOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	kafka.Consumer = consumer

	return kafka, nil
}

// newProducer initializes a new Kafka producer
func newProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// newConsumer initializes a new Kafka consumer
func newConsumer(brokers []string, group string, offset string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	if offset == "newest" {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// Close closes the Kafka producer and consumer
func (k *Kafka) Close() error {
	if err := k.Producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}
	return nil
}

// Produce sends a message to the Kafka producer topic
func (k *Kafka) Produce(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := k.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// Perform cleanup actions if necessary
	log.Println("Consumer group session ended")
	return nil
}

// ConsumeClaim starts a consumer loop of ConsumerGroupClaim's Messages()
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg.Value); err != nil {
			return err
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

// Consume reads messages from the Kafka consumer topic
func (k *Kafka) Consume(ctx context.Context, topic string, handler func([]byte) error) error {
	consumerHandler := &ConsumerGroupHandler{
		handler: handler,
		ready:   make(chan bool),
	}
	topics := []string{topic}

	go func() {
		for {
			if err := k.Consumer.Consume(ctx, topics, consumerHandler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			// Check if context is done
			if ctx.Err() != nil {
				return
			}
			consumerHandler.ready = make(chan bool)
		}
	}()

	// Await until the consumer has been set up
	<-consumerHandler.ready
	log.Println("Kafka consumer ready")
	return nil
}
