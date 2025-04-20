package repository

import (
	"context"

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

func (p *taskRepository) CreateTask(ctx context.Context, id string) error {
	return nil
}

func (p *taskRepository) GetTaskByID(ctx context.Context, id string) {

}
