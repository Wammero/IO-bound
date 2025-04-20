package main

import (
	"log"

	"github.com/Wammero/IO-bound/worker/internal/config"
	"github.com/Wammero/IO-bound/worker/internal/repository"
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.New(cfg.Database.GetConnStr())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()

	for true {

	}
}
