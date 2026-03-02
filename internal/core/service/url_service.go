package service

import (
	"context"
	"errors"
	"time"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"
	"urlshortener/internal/metrics"
	"urlshortener/internal/tracing"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
)

type urlService struct {
	urlRepo      out.URLRepository
	keyGenerator out.KeyGenerator
	hitRepo      out.URLHitEventRepository
	cache        out.URLCache
	validate     *validator.Validate
}

func NewUrlService(urlRepo out.URLRepository, keyGenerator out.KeyGenerator, hitRepo out.URLHitEventRepository, cache out.URLCache) in.URLService {
	return &urlService{
		urlRepo:      urlRepo,
		keyGenerator: keyGenerator,
		hitRepo:      hitRepo,
		cache:        cache,
		validate:     validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (u urlService) Create(ctx context.Context, originalURL string, customURL string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "urlService.Create")
	defer span.End()
	span.SetAttributes(attribute.String("url.original", originalURL))

	if err := u.validate.Var(originalURL, "required,url"); err != nil {
		logging.AppLogger.Debug("Invalid URL provided", "error", err.Error())
		tracing.RecordError(ctx, err)
		return "", err
	}

	if err := u.validate.Var(customURL, "omitempty,alphanum,min=3,max=20"); err != nil {
		logging.AppLogger.Debug("Invalid custom URL provided", "error", err.Error())
		tracing.RecordError(ctx, err)
		return "", err
	}

	if customURL == "" {
		for range 5 {
			customURL = u.keyGenerator.Generate()
			existing, err := u.urlRepo.FindByKey(ctx, customURL)
			if err != nil {
				logging.AppLogger.Error("Failed to check key existence", err)
				tracing.RecordError(ctx, err)
				return "", err
			}
			if existing == nil {
				break
			}
		}
	} else {
		existing, err := u.urlRepo.FindByKey(ctx, customURL)
		if err != nil {
			logging.AppLogger.Error("Failed to check custom key existence", err)
			tracing.RecordError(ctx, err)
			return "", err
		}
		if existing != nil {
			logging.AppLogger.Debug("Custom key already exists", "key", customURL)
			tracing.RecordError(ctx, errors.New("key already exists"))
			return "", errors.New("key already exists")
		}
	}

	url := &model.URL{
		ShortCode:   customURL,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	if err := u.urlRepo.Store(ctx, url); err != nil {
		logging.AppLogger.Error("Failed to store URL", err)
		tracing.RecordError(ctx, err)
		return "", err
	}
	span.SetAttributes(attribute.String("url.short_code", customURL))
	if err := u.cache.Set(ctx, customURL, originalURL); err != nil {
		logging.AppLogger.Warn("Failed to cache URL", "key", customURL, "error", err.Error())
	} else {
		if metrics.CacheSetsTotal != nil {
			metrics.CacheSetsTotal.Inc()
		}
	}

	if metrics.UrlsCreatedTotal != nil {
		metrics.UrlsCreatedTotal.Inc()
	}

	logging.AppLogger.Info("URL created successfully", "short_code", customURL, "original_url", originalURL)

	return customURL, nil
}

func (u urlService) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "urlService.GetOriginal")
	defer span.End()
	span.SetAttributes(attribute.String("url.short_code", shortURL))

	cached, err := u.cache.Get(ctx, shortURL)
	if err == nil {
		if cached != "" {
			if metrics.CacheHitsTotal != nil {
				metrics.CacheHitsTotal.Inc()
			}
		} else {
			if metrics.CacheMissesTotal != nil {
				metrics.CacheMissesTotal.Inc()
			}
		}
	}
	if err == nil && cached != "" {
		logging.AppLogger.Debug("Cache hit", "short_code", shortURL)
		span.SetAttributes(attribute.Bool("cache.hit", true))
		return cached, nil
	}

	logging.AppLogger.Debug("Cache miss", "short_code", shortURL)
	span.SetAttributes(attribute.Bool("cache.hit", false))

	url, err := u.urlRepo.FindByKey(ctx, shortURL)
	if err != nil {
		logging.AppLogger.Error("Failed to fetch URL from DB", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return "", err
	}
	if url == nil {
		logging.AppLogger.Debug("URL not found", "short_code", shortURL)
		tracing.RecordError(ctx, errors.New("URL not found"))
		return "", in.ErrNotFound
	}

	if err := u.cache.Set(ctx, shortURL, url.OriginalURL); err != nil {
		logging.AppLogger.Warn("Failed to cache URL after fetch", "key", shortURL, "error", err.Error())
	} else {
		if metrics.CacheSetsTotal != nil {
			metrics.CacheSetsTotal.Inc()
		}
	}

	logging.AppLogger.Debug("URL fetched from DB", "short_code", shortURL)

	return url.OriginalURL, nil
}

func (u urlService) GetAnalytics(ctx context.Context, shortURL string) (*model.Analytics, error) {
	ctx, span := tracing.StartSpan(ctx, "urlService.GetAnalytics")
	defer span.End()
	span.SetAttributes(attribute.String("url.short_code", shortURL))

	existing, err := u.urlRepo.FindByKey(ctx, shortURL)
	if err != nil {
		logging.AppLogger.Error("Failed to fetch URL for analytics", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}
	if existing == nil {
		logging.AppLogger.Debug("URL not found for analytics", "short_code", shortURL)
		tracing.RecordError(ctx, errors.New("URL not found"))
		return nil, errors.New("URL not found")
	}

	totalClicks, err := u.hitRepo.GetTotalClicks(ctx, shortURL)
	if err != nil {
		logging.AppLogger.Error("Failed to get total clicks", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	recentClicks, err := u.hitRepo.GetRecentClicks(ctx, shortURL, 100)
	if err != nil {
		logging.AppLogger.Error("Failed to get recent clicks", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	recentClickInfos := make([]model.ClickInfo, 0, len(recentClicks))
	for _, click := range recentClicks {
		recentClickInfos = append(recentClickInfos, model.ClickInfo{
			UserAgent: click.UserAgent,
			IP:        click.IP,
			Referer:   click.Referrer,
			Timestamp: click.Timestamp,
		})
	}

	byUserAgent, err := u.hitRepo.GetAggregatedByUserAgent(ctx, shortURL)
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by user agent", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	byDay, err := u.hitRepo.GetAggregatedByDay(ctx, shortURL, time.Time{}, time.Time{})
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by day", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	byMonth, err := u.hitRepo.GetAggregatedByMonth(ctx, shortURL, time.Time{}, time.Time{})
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by month", err, "short_code", shortURL)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	logging.AppLogger.Debug("Analytics fetched", "short_code", shortURL, "total_clicks", totalClicks)

	return &model.Analytics{
		ShortKey:     shortURL,
		TotalClicks:  totalClicks,
		RecentClicks: recentClickInfos,
		ByUserAgent:  byUserAgent,
		ByDay:        byDay,
		ByMonth:      byMonth,
	}, nil
}

func (u urlService) RecordHit(ctx context.Context, hit *model.URLHitEvent) error {
	ctx, span := tracing.StartSpan(ctx, "urlService.RecordHit")
	defer span.End()
	span.SetAttributes(attribute.String("hit.url_id", hit.URLID))

	if err := u.hitRepo.Store(ctx, hit); err != nil {
		logging.AppLogger.Error("Failed to record hit", err, "url_id", hit.URLID)
		tracing.RecordError(ctx, err)
		return err
	}
	logging.AppLogger.Debug("Hit recorded", "url_id", hit.URLID, "ip", hit.IP)
	return nil
}
