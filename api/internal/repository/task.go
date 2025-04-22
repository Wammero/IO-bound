package repository

import (
	"context"
	"fmt"

	"github.com/Wammero/IO-bound/api/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type taskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *taskRepository {
	return &taskRepository{pool: pool}
}

func (r *taskRepository) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *taskRepository) CreateTask(ctx context.Context, tx pgx.Tx, id, taskType, payload string) error {
	query := `
		INSERT INTO tasks (id, type, payload, status, created_at, updated_at)
		VALUES ($1, $2, $3::jsonb, 'pending', now(), now())
	`

	_, err := tx.Exec(ctx, query, id, taskType, payload)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	query := `SELECT id, type, payload, status, result, error, webhook_url, webhook_sent, created_at, updated_at
			  FROM tasks WHERE id = $1`

	var task models.Task

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Type,
		&task.Payload,
		&task.Status,
		&task.Result,
		&task.Error,
		&task.WebhookURL,
		&task.WebhookSent,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("task with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}
