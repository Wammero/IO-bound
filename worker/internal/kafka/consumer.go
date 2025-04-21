package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	noTimeout = 10 * time.Second
)

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
}

func NewConsumer(brokers, groupID, topic string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           groupID,
		"session.timeout.ms": 6000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Consume(ctx context.Context) {
	for {
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			log.Printf("Error while reading message: %v", err)
		}
		if kafkaMsg == nil {
			continue
		}
		if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg, TopicPartition.Offset); err != nil {
			log.Printf("Error while handling message: %v", err)
			continue
		}

	}
}

func (c *Consumer) Close() {
	log.Println("Closing Kafka consumer...")
	err := c.consumer.Close()
	if err != nil {
		log.Printf("Error while closing Kafka consumer: %v", err)
	}
}
