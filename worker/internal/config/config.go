package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Kafka    KafkaConfig
	Worker   WorkerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type KafkaConfig struct {
	Brokers string
	GroupID string
	Topics  KafkaTopics
}

type KafkaTopics struct {
	Topic1 string
	Topic2 string
	Topic3 string
}

type WorkerConfig struct {
	NumWorkers int
	LogLevel   string
}

func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnvOrFatal("DATABASE_HOST"),
			Port:     getEnvOrFatal("DATABASE_PORT"),
			User:     getEnvOrFatal("DATABASE_USER"),
			Password: getEnvOrFatal("DATABASE_PASSWORD"),
			Name:     getEnvOrFatal("DATABASE_NAME"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvOrFatal("KAFKA_BROKERS"),
			GroupID: getEnvOrFatal("KAFKA_GROUP_ID"),
			Topics: KafkaTopics{
				Topic1: getEnvOrFatal("KAFKA_TOPIC_1"),
				Topic2: getEnvOrFatal("KAFKA_TOPIC_2"),
				Topic3: getEnvOrFatal("KAFKA_TOPIC_3"),
			},
		},
		Worker: WorkerConfig{
			NumWorkers: getEnvIntOrFatal("WORKER_NUM_WORKERS"),
			LogLevel:   getEnvOrDefault("WORKER_LOG_LEVEL", "info"),
		},
	}
}

func (db DatabaseConfig) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name)
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return val
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvIntOrFatal(key string) int {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Environment variable %s should be an integer, got: %s", key, val)
	}
	return intVal
}
