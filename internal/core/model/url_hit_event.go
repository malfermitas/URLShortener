package model

import (
	"time"
)

type URLHitEvent struct {
	ID          int64     `json:"id"`
	URLID       string    `json:"url_id"`
	UserAgent   string    `json:"user_agent,omitempty"`
	IP          string    `json:"ip,omitempty"`
	CountryCode string    `json:"country_code,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	DeviceType  string    `json:"device_type,omitempty"`
	OS          string    `json:"os,omitempty"`
	Browser     string    `json:"browser,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}
