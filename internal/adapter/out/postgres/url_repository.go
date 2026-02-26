package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/out"
	"urlshortener/internal/logging"

	"github.com/wb-go/wbf/dbpg"
)

type urlRepository struct {
	db *dbpg.DB
}

func (u *urlRepository) GetDB() *dbpg.DB {
	return u.db
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
	query := `INSERT INTO urls (short_code, original_url, custom_code, created_at) VALUES ($1, $2, $3, $4)`
	_, err := u.db.ExecContext(ctx, query, url.ShortCode, url.OriginalURL, url.CustomCode, url.CreatedAt)
	if err != nil {
		logging.AppLogger.Error("Failed to store URL", err)
		return err
	}
	return nil
}

func (u urlRepository) FindByKey(ctx context.Context, shortKey string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, custom_code, created_at FROM urls WHERE short_code = $1`
	row := u.db.QueryRowContext(ctx, query, shortKey)

	var url model.URL
	err := row.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CustomCode, &url.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		logging.AppLogger.Error("Failed to find URL by key", err)
		return nil, err
	}
	return &url, nil
}
