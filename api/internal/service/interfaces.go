package service

import (
	"context"
)

type TaskService interface {
	CreateTask(ctx context.Context, id string) error
	GetTaskByID(ctx context.Context, id string)
}
