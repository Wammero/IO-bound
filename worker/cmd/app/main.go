package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Wammero/IO-bound/worker/internal/config"
	"github.com/Wammero/IO-bound/worker/internal/kafka"
	"github.com/Wammero/IO-bound/worker/internal/repository"
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.New(cfg.Database.GetConnStr())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем Ctrl+C
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Println("Shutting down consumer...")
		cancel()
	}()
	consumer.Consume(ctx)
}
