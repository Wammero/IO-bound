package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	TaskRepository TaskRepository
	pool           *pgxpool.Pool
}

func New(connstr string) (*Repository, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		TaskRepository: NewTaskRepository(pool),
		pool:           pool,
	}, nil
}

func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}
