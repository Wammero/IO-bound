package service

import (
	"context"

	"github.com/Wammero/IO-bound/api/internal/repository"
)

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *taskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, id string) error {
	return nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) {

}
