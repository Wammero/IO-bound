package main

import (
	"context"
	"log"

	"github.com/Wammero/IO-bound/worker/internal/config"
	"github.com/Wammero/IO-bound/worker/internal/repository"
	"github.com/Wammero/IO-bound/worker/internal/task_worker"
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.New(cfg.Database.GetConnStr())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()

	taskWorker := task_worker.NewTaskWorker(repo.TaskRepository)
	taskWorker.Start(context.Background(), cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topics.Topic1, cfg.Worker.NumWorkers, 5)
}
