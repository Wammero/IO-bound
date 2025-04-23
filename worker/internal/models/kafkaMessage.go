package models

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaMessage struct {
	Message  *kafka.Message
	Consumer *kafka.Consumer
}
