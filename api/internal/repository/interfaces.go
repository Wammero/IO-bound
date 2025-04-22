package repository

import (
	"context"

	"github.com/Wammero/IO-bound/api/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, tx pgx.Tx, id, taskType, payload string) error
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	Pool() *pgxpool.Pool
}
