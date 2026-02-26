package out

import (
	"context"
)

type URLCache interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Set(ctx context.Context, shortURL, originalURL string) error
	Delete(ctx context.Context, shortURL string) error
	Close() error
}
