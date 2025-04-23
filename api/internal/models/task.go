package models

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      string          `json:"status"`
	Result      json.RawMessage `json:"result"`
	Error       string          `json:"error"`
	WebhookURL  string          `json:"webhook_url"`
	WebhookSent bool            `json:"webhook_sent"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
