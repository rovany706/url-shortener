package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rovany706/url-shortener/internal/database"
)

var (
	insertEntrySQL = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url)
		VALUES ($1, $2)`, database.TableName)
	insertEntrySQLBatch = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url)
			VALUES ($1, $2)
			ON CONFLICT (full_url) DO NOTHING`, database.TableName)
	selectFullURLSQL = fmt.Sprintf(
		`SELECT full_url FROM %s
		WHERE short_id = $1`, database.TableName)
	selectShortIDSQL = fmt.Sprintf(
		`SELECT short_id FROM %s
		WHERE full_url = $1`, database.TableName)
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
	row := repository.db.DBConnection.QueryRowContext(ctx, selectFullURLSQL, shortID)
	err := row.Scan(&fullURL)

	if err != nil {
		return "", false
	}

	return fullURL, true
}

func (repository *DatabaseRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	_, err := repository.db.DBConnection.ExecContext(ctx, insertEntrySQL, shortID, fullURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
	}

	return err
}

func (repository *DatabaseRepository) GetShortID(ctx context.Context, fullURL string) (shortID string, err error) {
	row := repository.db.DBConnection.QueryRowContext(ctx, selectShortIDSQL, fullURL)
	err = row.Scan(&shortID)

	if err != nil {
		return "", err
	}

	return
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
		_, err := tx.ExecContext(ctx, insertEntrySQLBatch, shortID, fullURL)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
