package service

import "github.com/Wammero/IO-bound/api/internal/repository"

type Service struct {
	TaskService TaskService
}

func New(repo *repository.Repository) *Service {
	return &Service{
		TaskService: NewTaskService(repo.TaskRepository),
	}
}
