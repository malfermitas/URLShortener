package in

import (
	"context"
	"urlshortener/internal/core/model"
)

type URLService interface {
	// Create создаёт новую короткую ссылку.
	Create(ctx context.Context, originalURL string, customKey string) (string, error)

	// GetOriginal возвращает оригинальный URL по короткому ключу.
	GetOriginal(ctx context.Context, shortKey string) (string, error)

	// GetAnalytics возвращает аналитику по короткой ссылке.
	GetAnalytics(ctx context.Context, shortKey string) (*model.Analytics, error)
}
