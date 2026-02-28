package model

import (
	"time"
)

type URLHitEvent struct {
	ID          int64     `json:"id"`
	URLID       string    `json:"url_id"`
	UserAgent   string    `json:"user_agent,omitzero"`
	IP          string    `json:"ip,omitzero"`
	CountryCode string    `json:"country_code,omitzero"`
	Referrer    string    `json:"referrer,omitzero"`
	DeviceType  string    `json:"device_type,omitzero"`
	OS          string    `json:"os,omitzero"`
	Browser     string    `json:"browser,omitzero"`
	Timestamp   time.Time `json:"timestamp"`
}
