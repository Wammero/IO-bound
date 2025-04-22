package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer     *kafka.Producer
	deliveryWg   sync.WaitGroup
	deliveryChan chan *kafka.Message
	closeOnce    sync.Once
	ctx          context.Context
	cancel       context.CancelFunc
}

type Message struct {
	Topic string
	Key   string
	Value []byte
}

func New(brokers string, bufferSize int) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "all",
		"linger.ms":         5,
		"retries":           5,
		"retry.backoff.ms":  250,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	prod := &Producer{
		producer:     p,
		deliveryChan: make(chan *kafka.Message, bufferSize),
		ctx:          ctx,
		cancel:       cancel,
	}

	prod.deliveryWg.Add(1)
	go prod.handleDelivery()

	return prod, nil
}

func (p *Producer) Produce(msg Message) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.Topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(msg.Key),
		Value: msg.Value,
	}

	err := p.producer.Produce(kafkaMsg, nil)
	if err != nil {
		log.Printf("Produce error: %v", err)
		return err
	}

	return nil
}

func (p *Producer) handleDelivery() {
	defer p.deliveryWg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case ev := <-p.producer.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v", e.TopicPartition.Error)
				} else {
					log.Printf("Message delivered to %v [%d] at offset %v",
						*e.TopicPartition.Topic,
						e.TopicPartition.Partition,
						e.TopicPartition.Offset,
					)
				}
			default:
				log.Printf("Ignored event: %v", e)
			}
		}
	}
}

func (p *Producer) Close() {
	p.closeOnce.Do(func() {
		p.cancel()
		p.producer.Flush(5000)
		p.producer.Close()
		p.deliveryWg.Wait()
	})
}
