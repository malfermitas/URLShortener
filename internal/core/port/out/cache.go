package out

import (
	"context"
	"time"
)

// Cache предоставляет быстрый доступ к оригинальным URL по ключу.
type Cache interface {
	// Get возвращает оригинальный URL по ключу, если он есть в кэше.
	// Возвращает ErrCacheMiss, если ключ отсутствует.
	Get(ctx context.Context, key string) (string, error)

	// Set сохраняет оригинальный URL с временем жизни ttl.
	Set(ctx context.Context, key, value string, ttl time.Duration) error
}
