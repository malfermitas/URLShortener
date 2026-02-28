package model

import "time"

// Analytics — результат аналитики.
type Analytics struct {
	ShortKey     string           `json:"short_key"`
	TotalClicks  int64            `json:"total_clicks"`
	RecentClicks []ClickInfo      `json:"recent_clicks,omitzero"`
	ByUserAgent  map[string]int64 `json:"by_user_agent,omitzero"`
	ByDay        map[string]int64 `json:"by_day,omitzero"`
	ByMonth      map[string]int64 `json:"by_month,omitzero"`
}

type ClickInfo struct {
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Referer   string    `json:"referer"`
	Timestamp time.Time `json:"timestamp"`
}
