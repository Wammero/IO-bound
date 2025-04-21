package main

import (
	"log"
	"time"

	"github.com/Wammero/IO-bound/api/internal/config"
	"github.com/Wammero/IO-bound/api/internal/handler"
	"github.com/Wammero/IO-bound/api/internal/kafka"
	"github.com/Wammero/IO-bound/api/internal/migration"
	"github.com/Wammero/IO-bound/api/internal/repository"
	"github.com/Wammero/IO-bound/api/internal/router"
	"github.com/Wammero/IO-bound/api/internal/server"
	"github.com/Wammero/IO-bound/api/internal/service"
)

func main() {
	cfg := config.NewConfig()

	connstr := cfg.Database.GetConnStr()
	repo, err := repository.New(connstr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()

	migration.ApplyMigrations(connstr)

	kafkaProducer, err := kafka.New(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	svc := service.New(repo, kafkaProducer, cfg.Kafka.Topic)

	r := router.New()
	h := handler.New(svc)
	h.SetupRoutes(r)

	server.Start(server.Config{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
		Timeout: 5 * time.Second,
	})
}
