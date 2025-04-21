package service

import (
	"context"

	"github.com/Wammero/IO-bound/api/internal/kafka"
	"github.com/Wammero/IO-bound/api/internal/repository"
)

type taskService struct {
	repo     repository.TaskRepository
	producer *kafka.Producer
	topic    string
}

func NewTaskService(repo repository.TaskRepository, producer *kafka.Producer, topic string) *taskService {
	return &taskService{
		repo:     repo,
		producer: producer,
		topic:    topic,
	}
}

func (s *taskService) CreateTask(ctx context.Context, id string) error {
	message := "Привет"
	err := s.producer.SendMessage(s.topic, id, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) {

}
