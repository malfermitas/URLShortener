package out

import (
	"context"
	"urlshortener/internal/core/model"
)

// URLRepository отвечает за сохранение и поиск сокращённых ссылок.
type URLRepository interface {
	// Store сохраняет новую ссылку.
	Store(ctx context.Context, url *model.URL) error

	// FindByKey ищет ссылку по короткому ключу.
	// Возвращает nil, nil если не найдено.
	FindByKey(ctx context.Context, shortKey string) (*model.URL, error)

	// Exists проверяет, существует ли ключ.
	Exists(ctx context.Context, shortKey string) (bool, error)
}
