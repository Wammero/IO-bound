package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 15 * 1000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokers string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) SendMessage(topic string, key string, value []byte) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
		Key:            []byte(key),
	}

	deliveryChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return err
	}

	msg := <-deliveryChan
	switch ev := msg.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return fmt.Errorf("unknown event: %v", ev)
	}
}

func (p *Producer) Close() error {
	p.producer.Flush(flushTimeout)
	return p.producer.Close()
}
