package in

import (
	"context"
	"urlshortener/internal/core/model"
)

type URLService interface {
	// Create создаёт новую короткую ссылку.
	Create(ctx context.Context, originalURL string, customURL string) (string, error)

	// GetOriginal возвращает оригинальный URL по короткому ключу.
	GetOriginal(ctx context.Context, shortURL string) (string, error)

	// GetAnalytics возвращает аналитику по короткой ссылке.
	GetAnalytics(ctx context.Context, shortURL string) (*model.Analytics, error)
}
