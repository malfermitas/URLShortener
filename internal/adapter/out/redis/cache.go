package redis

import (
	"context"
	"errors"
	"fmt"
	"time"
	"urlshortener/internal/adapter/out/retry"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"
	"urlshortener/internal/tracing"

	"github.com/wb-go/wbf/redis"
)

type urlCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewURLCache(addr string, password string, db int, maxTTL time.Duration) (out.URLCache, error) {
	client := redis.New(addr, password, db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		logging.AppLogger.Error("Cannot connect to Redis", err)
		return nil, err
	}

	logging.AppLogger.Info("Connected to Redis")
	return &urlCache{
		client: client,
		ttl:    maxTTL,
	}, nil
}

func (c *urlCache) Get(ctx context.Context, shortURL string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "redis.Get")
	defer span.End()

	key := fmt.Sprintf("url:%s", shortURL)
	strategy := retry.GetRedisStrategy()

	data, err := c.client.GetWithRetry(ctx, strategy, key)
	if errors.Is(err, redis.NoMatches) {
		return "", nil
	}
	if err != nil {
		logging.AppLogger.Error("Failed to get URL from cache", err)
		tracing.RecordError(ctx, err)
		return "", err
	}

	return data, nil
}

func (c *urlCache) Set(ctx context.Context, shortURL, originalURL string) error {
	ctx, span := tracing.StartSpan(ctx, "redis.Set")
	defer span.End()

	key := fmt.Sprintf("url:%s", shortURL)
	strategy := retry.GetRedisStrategy()

	if err := c.client.SetWithExpirationAndRetry(ctx, strategy, key, originalURL, c.ttl); err != nil {
		logging.AppLogger.Error("Failed to set URL in cache", err)
		tracing.RecordError(ctx, err)
		return err
	}

	return nil
}

func (c *urlCache) Delete(ctx context.Context, shortKey string) error {
	ctx, span := tracing.StartSpan(ctx, "redis.Delete")
	defer span.End()

	key := fmt.Sprintf("url:%s", shortKey)
	strategy := retry.GetRedisStrategy()

	if err := c.client.DelWithRetry(ctx, strategy, key); err != nil {
		logging.AppLogger.Error("Failed to delete URL from cache", err)
		tracing.RecordError(ctx, err)
		return err
	}
	return nil
}

func (c *urlCache) Close() error {
	return c.client.Close()
}
