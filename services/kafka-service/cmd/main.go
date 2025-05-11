package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Get Kafka brokers from environment or use default
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(kafkaBrokers, ",")

	// Create a new Kafka consumer config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0 // Use appropriate Kafka version

	// Create a context that will be canceled on SIGTERM or SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create a wait group to manage goroutines
	var wg sync.WaitGroup

	// Start the consumer in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeOrders(ctx, brokers, config)
	}()

	// Start the producer in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		produceOrderUpdates(ctx, brokers, config)
	}()

	// Wait for termination signal
	<-signals
	log.Println("Received termination signal. Shutting down...")
	cancel()

	// Wait for goroutines to finish
	wg.Wait()
	log.Println("Kafka Order Service shut down successfully")
}

// consumeOrders consumes messages from the orders topic
func consumeOrders(ctx context.Context, brokers []string, config *sarama.Config) {
	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "order-service", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Create a consumer handler
	handler := orderConsumerHandler{}

	// Consume messages in a loop
	for {
		// Check if context is canceled (service is shutting down)
		select {
		case <-ctx.Done():
			log.Println("Stopping order consumer...")
			return
		default:
			// Continue consuming
		}

		// Consume messages from the orders topic
		err := consumerGroup.Consume(ctx, []string{"orders"}, handler)
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}
}

// produceOrderUpdates simulates producing order status updates
func produceOrderUpdates(ctx context.Context, brokers []string, config *sarama.Config) {
	// Create a new sync producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	// Simulate producing order updates every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping order producer...")
			return
		case <-ticker.C:
			// Create a mock order update
			message := &sarama.ProducerMessage{
				Topic: "order_updates",
				Key:   sarama.StringEncoder("order-123"),
				Value: sarama.StringEncoder(`{"id":"order-123","status":"PROCESSING"}`),
			}

			// Send the message
			_, _, err := producer.SendMessage(message)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			} else {
				log.Println("Order update sent successfully")
			}
		}
	}
}

// orderConsumerHandler implements sarama.ConsumerGroupHandler interface
type orderConsumerHandler struct{}

func (h orderConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h orderConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h orderConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received order: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		
		// Process the order (add your business logic here)
		
		// Mark the message as processed
		session.MarkMessage(message, "")
	}
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
