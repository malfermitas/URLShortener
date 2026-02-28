package in

import (
	"context"
	"errors"
	"urlshortener/internal/core/model"
)

var (
	ErrNotFound = errors.New("URL not found")
)

type URLService interface {
	// Create создаёт новую короткую ссылку.
	Create(ctx context.Context, originalURL string, customURL string) (string, error)

	// GetOriginal возвращает оригинальный URL по короткому ключу.
	GetOriginal(ctx context.Context, shortURL string) (string, error)

	// GetAnalytics возвращает аналитику по короткой ссылке.
	GetAnalytics(ctx context.Context, shortURL string) (*model.Analytics, error)

	// RecordHit записывает факт перехода по ссылке.
	RecordHit(ctx context.Context, hit *model.URLHitEvent) error
}
