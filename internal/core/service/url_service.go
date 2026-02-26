package service

import (
	"context"
	"errors"
	"time"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/core/port/out"
)

type urlService struct {
	urlRepo      out.URLRepository
	keyGenerator out.KeyGenerator
	hitRepo      out.URLHitEventRepository
	cache        out.URLCache
}

func NewUrlService(urlRepo out.URLRepository, keyGenerator out.KeyGenerator, hitRepo out.URLHitEventRepository, cache out.URLCache) in.URLService {
	return &urlService{
		urlRepo:      urlRepo,
		keyGenerator: keyGenerator,
		hitRepo:      hitRepo,
		cache:        cache,
	}
}

func (u urlService) Create(ctx context.Context, originalURL string, customURL string) (string, error) {
	shortKey := customURL
	if shortKey == "" {
		for {
			shortKey = u.keyGenerator.Generate()
			existing, err := u.urlRepo.FindByKey(ctx, shortKey)
			if err != nil {
				return "", err
			}
			if existing == nil {
				break
			}
		}
	} else {
		existing, err := u.urlRepo.FindByKey(ctx, shortKey)
		if err != nil {
			return "", err
		}
		if existing != nil {
			return "", errors.New("key already exists")
		}
	}

	url := &model.URL{
		ShortCode:   shortKey,
		OriginalURL: originalURL,
		CustomCode:  customURL,
		CreatedAt:   time.Now(),
	}

	if err := u.urlRepo.Store(ctx, url); err != nil {
		return "", err
	}

	_ = u.cache.Set(ctx, customURL, originalURL)

	return shortKey, nil
}

func (u urlService) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	cached, err := u.cache.Get(ctx, shortURL)
	if err == nil {
		return cached, nil
	}

	url, err := u.urlRepo.FindByKey(ctx, shortURL)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", errors.New("URL not found")
	}

	_ = u.cache.Set(ctx, shortURL, url.OriginalURL)

	return url.OriginalURL, nil
}

func (u urlService) GetAnalytics(ctx context.Context, shortURL string) (*model.Analytics, error) {
	existing, err := u.urlRepo.FindByKey(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("URL not found")
	}

	totalClicks, err := u.hitRepo.GetTotalClicks(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	recentClicks, err := u.hitRepo.GetRecentClicks(ctx, shortURL, 100)
	if err != nil {
		return nil, err
	}

	recentClickInfos := make([]model.ClickInfo, len(recentClicks))
	for i, click := range recentClicks {
		recentClickInfos[i] = model.ClickInfo{
			UserAgent: click.UserAgent,
			IP:        click.IP,
			Referer:   click.Referrer,
			Timestamp: click.Timestamp,
		}
	}

	byUserAgent, err := u.hitRepo.GetAggregatedByUserAgent(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	byDay, err := u.hitRepo.GetAggregatedByDay(ctx, shortURL, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	byMonth, err := u.hitRepo.GetAggregatedByMonth(ctx, shortURL, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	return &model.Analytics{
		ShortKey:     shortURL,
		TotalClicks:  totalClicks,
		RecentClicks: recentClickInfos,
		ByUserAgent:  byUserAgent,
		ByDay:        byDay,
		ByMonth:      byMonth,
	}, nil
}
