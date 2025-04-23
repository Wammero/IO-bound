package task_worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Wammero/IO-bound/worker/internal/kafka"
	"github.com/Wammero/IO-bound/worker/internal/repository"
)

type TaskWorker struct {
	repo repository.TaskRepository
}

func NewTaskWorker(repo repository.TaskRepository) *TaskWorker {
	return &TaskWorker{
		repo: repo,
	}
}

func (w *TaskWorker) Start(ctx context.Context, brokers string, groupID string, topic string, numConsumers int, numGoroutines int) {
	messageChannel, err := kafka.StartConsumers(ctx, brokers, groupID, topic, numConsumers)
	if err != nil {
		log.Fatalf("Failed to create consumers: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range messageChannel {
				checked, err := w.repo.TaskAlreadyChecked(ctx, string(msg.Message.Key))
				if err != nil {
					log.Printf("Failed to check task existence: %v", err)
					w.repo.UpdateTaskStatus(ctx, string(msg.Message.Key), "failed")
					continue
				}
				if !checked {
					log.Printf("Received message: %s", string(msg.Message.Value))
					w.repo.UpdateTaskStatus(ctx, string(msg.Message.Key), "in_progress")
					time.Sleep(20 * time.Second)
					// if err != nil {
					// 	w.repo.UpdateTaskStatus(ctx, string(msg.Key), "failed")
					// }
					w.repo.UpdateTaskStatus(ctx, string(msg.Message.Key), "done")
					if err != nil {
						log.Printf("Failed to commit message: %v", err)
					}
					_, err = msg.Consumer.CommitMessage(msg.Message)
					if err != nil {
						log.Printf("Failed to commit Kafka message: %v", err)
					}
				}
				// else // just skip

			}
		}()
	}

	wg.Wait()
}
