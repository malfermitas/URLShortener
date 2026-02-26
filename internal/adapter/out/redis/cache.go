package redis

import (
	"context"
	"errors"
	"fmt"
	"time"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"

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
	key := fmt.Sprintf("url:%s", shortURL)
	data, err := c.client.Get(ctx, key)
	if errors.Is(err, redis.NoMatches) {
		return "", nil
	}
	if err != nil {
		logging.AppLogger.Error("Failed to get URL from cache", err)
		return "", err
	}

	return data, nil
}

func (c *urlCache) Set(ctx context.Context, shortURL, originalURL string) error {
	key := fmt.Sprintf("url:%s", shortURL)

	if err := c.client.SetWithExpiration(ctx, key, originalURL, c.ttl); err != nil {
		logging.AppLogger.Error("Failed to set URL in cache", err)
		return err
	}

	return nil
}

func (c *urlCache) Delete(ctx context.Context, shortKey string) error {
	key := fmt.Sprintf("url:%s", shortKey)
	if err := c.client.Del(ctx, key); err != nil {
		logging.AppLogger.Error("Failed to delete URL from cache", err)
		return err
	}
	return nil
}

func (c *urlCache) Close() error {
	return c.client.Close()
}
