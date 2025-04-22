package repository

import (
	"context"
	"fmt"

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

func (r *taskRepository) UpdateTaskStatus(ctx context.Context, taskID string, status string) error {
	query := `UPDATE tasks SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, taskID)
	return err
}

func (r *taskRepository) TaskAlreadyChecked(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1 AND status <> 'pending')`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check task existence: %v", err)
	}

	return exists, nil
}
