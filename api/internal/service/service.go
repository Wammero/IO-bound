package service

import (
	"github.com/Wammero/IO-bound/api/internal/config"
	"github.com/Wammero/IO-bound/api/internal/kafka"
	"github.com/Wammero/IO-bound/api/internal/repository"
)

type Service struct {
	TaskService TaskService
}

func New(repo *repository.Repository, producer *kafka.Producer, topics config.KafkaTopics) *Service {
	return &Service{
		TaskService: NewTaskService(repo.TaskRepository, producer, topics),
	}
}
