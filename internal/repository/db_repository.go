package repository

import (
	"context"
	"fmt"

	"github.com/rovany706/url-shortener/internal/database"
)

var (
	insertEntrySql = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url)
		VALUES ($1, $2)`, database.TableName)
	selectEntrySql = fmt.Sprintf(
		`SELECT full_url FROM %s
		WHERE short_id = $1 LIMIT 1`, database.TableName)
)

type DBRepository struct {
	db database.Database
}

func NewDBRepository(ctx context.Context, connString string) (Repository, error) {
	db, err := database.InitConnection(ctx, connString)

	if err != nil {
		return nil, err
	}

	dbRepository := DBRepository{db: db}

	if err = dbRepository.db.EnsureCreated(ctx); err != nil {
		return nil, err
	}

	return &dbRepository, nil
}

func (repository *DBRepository) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	row := repository.db.QueryRowContext(ctx, selectEntrySql, shortID)
	err := row.Scan(&fullURL)

	if err != nil {
		return "", false
	}

	return fullURL, true
}

func (repository *DBRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	_, err := repository.db.ExecContext(ctx, insertEntrySql, shortID, fullURL)

	return err
}

func (repository *DBRepository) Close() {
	repository.db.Close()
}

func (repository *DBRepository) Ping(ctx context.Context) error {
	return repository.db.PingContext(ctx)
}
