package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Kafka    KafkaConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port string
}

type KafkaConfig struct {
	Brokers string
	Topics  KafkaTopics
	GroupID string
}

type KafkaTopics struct {
	Topic1 string
	Topic2 string
	Topic3 string
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
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvOrFatal("KAFKA_BROKERS"),
			Topics: KafkaTopics{
				Topic1: getEnvOrFatal("KAFKA_TOPIC_1"),
				Topic2: getEnvOrFatal("KAFKA_TOPIC_2"),
				Topic3: getEnvOrFatal("KAFKA_TOPIC_3"),
			},
			GroupID: getEnvOrFatal("KAFKA_GROUP_ID"),
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

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
