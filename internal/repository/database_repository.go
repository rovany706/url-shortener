package repository

import (
	"context"
	"fmt"

	"github.com/rovany706/url-shortener/internal/database"
)

var (
	insertEntrySQL = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url)
		VALUES ($1, $2)
		ON CONFLICT (short_id) DO NOTHING`, database.TableName)
	selectEntrySQL = fmt.Sprintf(
		`SELECT full_url FROM %s
		WHERE short_id = $1 LIMIT 1`, database.TableName)
)

type DatabaseRepository struct {
	db *database.Database
}

func NewDatabaseRepository(ctx context.Context, connString string) (Repository, error) {
	db, err := database.InitConnection(ctx, connString)

	if err != nil {
		return nil, err
	}

	dbRepository := DatabaseRepository{db: db}

	if err = dbRepository.db.EnsureCreated(ctx); err != nil {
		return nil, err
	}

	return &dbRepository, nil
}

func (repository *DatabaseRepository) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	row := repository.db.DBConnection.QueryRowContext(ctx, selectEntrySQL, shortID)
	err := row.Scan(&fullURL)

	if err != nil {
		return "", false
	}

	return fullURL, true
}

func (repository *DatabaseRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	_, err := repository.db.DBConnection.ExecContext(ctx, insertEntrySQL, shortID, fullURL)

	return err
}

func (repository *DatabaseRepository) Close() error {
	return repository.db.DBConnection.Close()
}

func (repository *DatabaseRepository) Ping(ctx context.Context) error {
	return repository.db.DBConnection.PingContext(ctx)
}

func (repository *DatabaseRepository) SaveEntries(ctx context.Context, shortIDMap map[string]string) error {
	tx, err := repository.db.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for shortID, fullURL := range shortIDMap {
		_, err := tx.ExecContext(ctx, insertEntrySQL, shortID, fullURL)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
