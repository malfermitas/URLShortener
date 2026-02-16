package model

import "time"

// Analytics — результат аналитики.
type Analytics struct {
	ShortKey     string
	TotalClicks  int64
	RecentClicks []ClickInfo      // последние N переходов (например, 100)
	ByUserAgent  map[string]int64 // количество переходов по user-agent
	ByDay        map[string]int64 // агрегация по дням (ключ: YYYY-MM-DD)
	ByMonth      map[string]int64 // агрегация по месяцам (ключ: YYYY-MM)
}

// ClickInfo — сокращённая информация о клике для аналитики.
type ClickInfo struct {
	UserAgent string
	IP        string
	Referer   string
	Timestamp time.Time
}
