package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	localretry "urlshortener/internal/adapter/out/retry"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"
	"urlshortener/internal/tracing"

	"github.com/wb-go/wbf/dbpg"
)

type urlRepository struct {
	db *dbpg.DB
}

func NewURLRepository(dsn string) (out.URLRepository, error) {
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

	if err = db.Master.Ping(); err != nil {
		logging.AppLogger.Error("Cannot connect to postgres database", err)
		return nil, err
	}

	return &urlRepository{
		db: db,
	}, nil
}

func (u urlRepository) Store(ctx context.Context, url *model.URL) error {
	ctx, span := tracing.StartSpan(ctx, "postgres.Store")
	defer span.End()

	query := `INSERT INTO urls (short_code, original_url, created_at) VALUES ($1, $2, $3)`
	strategy := localretry.GetDatabaseStrategy()
	_, err := u.db.ExecWithRetry(ctx, strategy, query, url.ShortCode, url.OriginalURL, url.CreatedAt)
	if err != nil {
		logging.AppLogger.Error("Failed to store URL", err)
		tracing.RecordError(ctx, err)
		return err
	}
	return nil
}

func (u urlRepository) FindByKey(ctx context.Context, shortKey string) (*model.URL, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.FindByKey")
	defer span.End()

	query := `SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = $1`
	row, err := u.db.QueryRowWithRetry(ctx, localretry.GetDatabaseStrategy(), query, shortKey)
	if err != nil {
		logging.AppLogger.Error("Failed to fetch URL", err)
		tracing.RecordError(ctx, err)
		return nil, err
	}

	var url model.URL
	err = row.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		logging.AppLogger.Error("Failed to find URL by key", err)
		tracing.RecordError(ctx, err)
		return nil, err
	}
	return &url, nil
}
