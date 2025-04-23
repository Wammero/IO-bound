package service

import (
	"context"

	"github.com/Wammero/IO-bound/api/internal/models"
)

type TaskService interface {
	CreateTask(ctx context.Context, id, taskType, payload string) (string, error)
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
}
