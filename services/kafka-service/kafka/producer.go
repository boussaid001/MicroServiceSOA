package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true // This is required for SyncProducer

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// SendMessage sends a message to Kafka
func (p *Producer) SendMessage(key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
