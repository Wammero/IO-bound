package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 15 * 1000 // in milliseconds
)

type Producer struct {
	producer *kafka.Producer
}

// New создает нового Kafka-продюсера
func New(brokers string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers, // Исправлен ключ
		"acks":              "all",
		"linger.ms":         5,
		"retries":           3,
		"retry.backoff.ms":  250,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

// SendMessage отправляет сообщение в указанный топик Kafka
func (p *Producer) SendMessage(topic string, key string, value []byte) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
		Key:   []byte(key),
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	if err := p.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case ev := <-deliveryChan:
		switch msg := ev.(type) {
		case *kafka.Message:
			if msg.TopicPartition.Error != nil {
				log.Printf("delivery failed: %v", msg.TopicPartition.Error)
				return msg.TopicPartition.Error
			}
			log.Printf("message delivered to topic %s [%d] at offset %v",
				*msg.TopicPartition.Topic,
				msg.TopicPartition.Partition,
				msg.TopicPartition.Offset,
			)
			return nil
		case kafka.Error:
			log.Printf("kafka error event: %v", msg)
			return msg
		default:
			return fmt.Errorf("unknown event type received: %T", ev)
		}
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for delivery report")
	}
}

// Close завершает работу продюсера и очищает буферы
func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
