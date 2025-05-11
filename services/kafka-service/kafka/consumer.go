package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// MessageHandler is a function that processes Kafka messages
type MessageHandler func(message *sarama.ConsumerMessage) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  MessageHandler
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
	}, nil
}

// Consume starts consuming messages from Kafka
func (c *Consumer) Consume(ctx context.Context) error {
	// Create a new consumer handler
	handler := consumerGroupHandler{
		handler: c.handler,
	}

	// Start consuming in a loop
	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Consume messages
		err := c.consumer.Consume(ctx, c.topics, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler implements the sarama.ConsumerGroupHandler interface
type consumerGroupHandler struct {
	handler MessageHandler
}

func (h consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := h.handler(message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		} else {
			session.MarkMessage(message, "")
		}
	}
	return nil
}
