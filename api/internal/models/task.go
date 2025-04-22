package models

type Task struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Payload     string `json:"payload"`
	Status      string `json:"status"`
	Result      string `json:"result"`
	Error       string `json:"error"`
	WebhookURL  string `json:"webhook_url"`
	WebhookSent bool   `json:"webhook_sent"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
