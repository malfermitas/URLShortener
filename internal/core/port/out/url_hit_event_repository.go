package out

import (
	"context"
	"time"
	"urlshortener/internal/core/model"
)

// URLHitEventRepository отвечает за запись и чтение статистики переходов.
type URLHitEventRepository interface {
	// Store сохраняет информацию о переходе.
	Store(ctx context.Context, click *model.URLHitEvent) error

	// GetTotalClicks возвращает общее количество переходов по ключу.
	GetTotalClicks(ctx context.Context, shortKey string) (int64, error)

	// GetRecentClicks возвращает последние limit переходов.
	GetRecentClicks(ctx context.Context, shortKey string, limit int) ([]model.URLHitEvent, error)

	// GetAggregatedByUserAgent возвращает количество переходов по user-agent.
	GetAggregatedByUserAgent(ctx context.Context, shortKey string) (map[string]int64, error)

	// GetAggregatedByDay возвращает количество переходов по дням за указанный период.
	// from и to могут быть нулевыми для периода "всё время".
	GetAggregatedByDay(ctx context.Context, shortKey string, from, to time.Time) (map[string]int64, error)

	// GetAggregatedByMonth аналогично по месяцам.
	GetAggregatedByMonth(ctx context.Context, shortKey string, from, to time.Time) (map[string]int64, error)
}
