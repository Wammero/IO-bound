package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskRepository interface {
	UpdateTaskStatus(ctx context.Context, taskID string, status string) error
	TaskAlreadyChecked(ctx context.Context, id string) (bool, error)
	Pool() *pgxpool.Pool
}
