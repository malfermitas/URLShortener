package postgres

import (
	"context"
	"fmt"
	"time"
	"urlshortener/internal/adapter/out/retry"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"

	"github.com/wb-go/wbf/dbpg"
)

type urlHitEventRepository struct {
	db *dbpg.DB
}

func NewURLHitEventRepository(dsn string) (out.URLHitEventRepository, error) {
	dbOptions := &dbpg.Options{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := dbpg.New(dsn, nil, dbOptions)
	if err != nil {
		logging.AppLogger.Error("Cannot connect to postgres database", err)
		return nil, err
	}

	return &urlHitEventRepository{
		db: db,
	}, nil
}

func (u urlHitEventRepository) Store(ctx context.Context, click *model.URLHitEvent) error {
	query := `INSERT INTO url_hit_events (url_id, user_agent, ip, country_code, referrer, device_type, os, browser, timestamp)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	strategy := retry.GetDatabaseStrategy()
	_, err := u.db.ExecWithRetry(
		ctx,
		strategy,
		query,
		click.URLID,
		click.UserAgent,
		click.IP,
		click.CountryCode,
		click.Referrer,
		click.DeviceType,
		click.OS,
		click.Browser,
		click.Timestamp,
	)
	if err != nil {
		logging.AppLogger.Error("Failed to store click event", err)
		return err
	}
	return nil
}

func (u urlHitEventRepository) GetTotalClicks(ctx context.Context, shortKey string) (int64, error) {
	query := `SELECT COUNT(*) FROM url_hit_events WHERE url_id = $1`
	strategy := retry.GetDatabaseStrategy()
	row, err := u.db.QueryRowWithRetry(ctx, strategy, query, shortKey)
	if err != nil {
		logging.AppLogger.Error("Failed to get total clicks", err)
		return 0, err
	}
	var count int64

	err = row.Scan(&count)
	if err != nil {
		logging.AppLogger.Error("Failed to get total clicks", err)
		return 0, err
	}
	return count, nil
}

func (u urlHitEventRepository) GetRecentClicks(ctx context.Context, shortKey string, limit int) ([]model.URLHitEvent, error) {
	query := `SELECT id, url_id, user_agent, ip, country_code, referrer, device_type, os, browser, timestamp
			  FROM url_hit_events WHERE url_id = $1 ORDER BY timestamp DESC LIMIT $2`
	rows, err := u.db.QueryContext(ctx, query, shortKey, limit)
	if err != nil {
		logging.AppLogger.Error("Failed to get recent clicks", err)
		return nil, err
	}
	defer rows.Close()

	var clicks []model.URLHitEvent
	for rows.Next() {
		var click model.URLHitEvent
		err := rows.Scan(&click.ID, &click.URLID, &click.UserAgent, &click.IP, &click.CountryCode, &click.Referrer, &click.DeviceType, &click.OS, &click.Browser, &click.Timestamp)
		if err != nil {
			logging.AppLogger.Error("Failed to scan click", err)
			return nil, err
		}
		clicks = append(clicks, click)
	}
	return clicks, nil
}

func (u urlHitEventRepository) GetAggregatedByUserAgent(ctx context.Context, shortKey string) (map[string]int64, error) {
	query := `SELECT user_agent, COUNT(*) as cnt FROM url_hit_events WHERE url_id = $1 GROUP BY user_agent`
	rows, err := u.db.QueryContext(ctx, query, shortKey)
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by user agent", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var userAgent string
		var count int64
		if err := rows.Scan(&userAgent, &count); err != nil {
			return nil, err
		}
		result[userAgent] = count
	}
	return result, nil
}

func (u urlHitEventRepository) GetAggregatedByDay(ctx context.Context, shortKey string, from, to time.Time) (map[string]int64, error) {
	query := `SELECT DATE(timestamp)::text as day, COUNT(*) as cnt FROM url_hit_events WHERE url_id = $1`
	args := []any{shortKey}

	if !from.IsZero() {
		args = append(args, from)
		query += fmt.Sprintf(" AND timestamp >= $%d", len(args))
	}
	if !to.IsZero() {
		args = append(args, to)
		query += fmt.Sprintf(" AND timestamp <= $%d", len(args))
	}

	query += " GROUP BY day ORDER BY day"

	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by day", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var day string
		var count int64
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result[day] = count
	}
	return result, nil
}

func (u urlHitEventRepository) GetAggregatedByMonth(ctx context.Context, shortKey string, from, to time.Time) (map[string]int64, error) {
	query := `SELECT TO_CHAR(timestamp, 'YYYY-MM') as month, COUNT(*) as cnt FROM url_hit_events WHERE url_id = $1`
	args := []any{shortKey}

	if !from.IsZero() {
		args = append(args, from)
		query += fmt.Sprintf(" AND timestamp >= $%d", len(args))
	}
	if !to.IsZero() {
		args = append(args, to)
		query += fmt.Sprintf(" AND timestamp <= $%d", len(args))
	}

	query += " GROUP BY month ORDER BY month"

	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		logging.AppLogger.Error("Failed to get aggregated by month", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var month string
		var count int64
		if err := rows.Scan(&month, &count); err != nil {
			return nil, err
		}
		result[month] = count
	}
	return result, nil
}
