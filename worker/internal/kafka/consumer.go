package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	noTimeout = 10 * time.Second
)

type Consumer struct {
	consumer *kafka.Consumer
}

// NewConsumer создает несколько Kafka потребителей.
func NewConsumer(brokers, groupID, topic string, numConsumers int) ([]*Consumer, error) {
	var consumers []*Consumer

	for i := 0; i < numConsumers; i++ {
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":               brokers,
			"group.id":                        groupID,
			"auto.offset.reset":               "earliest",
			"enable.auto.commit":              true,
			"auto.commit.interval.ms":         5000,
			"session.timeout.ms":              6000,
			"heartbeat.interval.ms":           2000,
			"max.poll.interval.ms":            300000,
			"go.events.channel.enable":        true,
			"go.application.rebalance.enable": true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create consumer: %w", err)
		}

		if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
			return nil, fmt.Errorf("failed to subscribe to topics: %w", err)
		}

		consumers = append(consumers, &Consumer{consumer: c})
	}

	return consumers, nil
}

func (c *Consumer) Consume(ctx context.Context, messageChannel chan<- *kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Starting to consume messages from Kafka...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping consumer")
			return
		default:
			msg, err := c.consumer.ReadMessage(noTimeout)
			if err != nil {
				var kafkaErr kafka.Error
				if ok := errors.As(err, &kafkaErr); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("Error while reading message: %v", err)
				continue
			}

			if msg != nil {
				log.Printf("Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s",
					*msg.TopicPartition.Topic,
					msg.TopicPartition.Partition,
					msg.TopicPartition.Offset,
					string(msg.Key),
					string(msg.Value))

				messageChannel <- msg
			}
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

func StartConsumers(ctx context.Context, brokers, groupID, topic string, numConsumers int) (chan *kafka.Message, error) {
	consumers, err := NewConsumer(brokers, groupID, topic, numConsumers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create consumers: %v", err)
	}

	messageChannel := make(chan *kafka.Message, 100)

	var wg sync.WaitGroup
	for _, consumer := range consumers {
		wg.Add(1)
		go consumer.Consume(ctx, messageChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(messageChannel)
	}()

	return messageChannel, nil
}
