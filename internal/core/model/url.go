package model

import (
	"time"
)

type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CustomCode  string    `json:"custom_code,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
