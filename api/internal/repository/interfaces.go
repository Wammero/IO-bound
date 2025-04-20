package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, id string) error
	GetTaskByID(ctx context.Context, id string)
	Pool() *pgxpool.Pool
}
