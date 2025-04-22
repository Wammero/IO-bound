package service

import (
	"context"
	"fmt"

	"github.com/Wammero/IO-bound/api/internal/config"
	"github.com/Wammero/IO-bound/api/internal/kafka"
	"github.com/Wammero/IO-bound/api/internal/models"
	"github.com/Wammero/IO-bound/api/internal/repository"
)

type taskService struct {
	repo     repository.TaskRepository
	producer *kafka.Producer
	topics   config.KafkaTopics
}

func NewTaskService(repo repository.TaskRepository, producer *kafka.Producer, topics config.KafkaTopics) *taskService {
	return &taskService{
		repo:     repo,
		producer: producer,
		topics:   topics,
	}
}

func (s *taskService) CreateTask(ctx context.Context, id, taskType, payload string) error {
	tx, err := s.repo.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()
	err = s.repo.CreateTask(ctx, tx, id, taskType, payload)
	if err != nil {
		return fmt.Errorf("не удалось создать задачу: %v", err)
	}
	message := id
	err = s.producer.Produce(kafka.Message{
		Topic: s.topics.Topic1,
		Key:   id,
		Value: []byte(message),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}
