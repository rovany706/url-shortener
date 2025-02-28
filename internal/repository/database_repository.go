package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rovany706/url-shortener/internal/database"
	"github.com/rovany706/url-shortener/internal/models"
)

var (
	insertEntrySQL = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url, user_id, is_deleted)
		VALUES ($1, $2, $3, false)`, database.ShortLinksTableName)
	insertEntrySQLBatch = fmt.Sprintf(
		`INSERT INTO %s (short_id, full_url, user_id, is_deleted)
			VALUES ($1, $2, $3, false)
			ON CONFLICT (full_url) DO NOTHING`, database.ShortLinksTableName)
	selectFullURLSQL = fmt.Sprintf(
		`SELECT full_url FROM %s
		WHERE short_id = $1`, database.ShortLinksTableName)
	selectShortIDSQL = fmt.Sprintf(
		`SELECT short_id FROM %s
		WHERE full_url = $1`, database.ShortLinksTableName)
	selectUserURLs = fmt.Sprintf(
		`SELECT short_id, full_url FROM %s
		WHERE user_id = $1`, database.ShortLinksTableName)
	insertNewUser = fmt.Sprintf(
		`INSERT INTO %s DEFAULT VALUES RETURNING id;`,
		database.UsersTableName)
	deleteShortLink = fmt.Sprintf(
		`UPDATE %s
		SET is_deleted = true
		WHERE short_id = $1 AND user_id = $2`,
		database.ShortLinksTableName)
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

func (repository *DatabaseRepository) SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error {
	_, err := repository.db.DBConnection.ExecContext(ctx, insertEntrySQL, shortID, fullURL, userID)

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

func (repository *DatabaseRepository) GetUserEntries(ctx context.Context, userID int) (shortIDMap ShortIDMap, err error) {
	rows, err := repository.db.DBConnection.QueryContext(ctx, selectUserURLs, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shortIDMap = make(ShortIDMap, 0)
	for rows.Next() {
		var userEntry struct {
			shortID string
			fullURL string
		}

		err = rows.Scan(&userEntry.shortID, &userEntry.fullURL)
		if err != nil {
			return nil, err
		}

		shortIDMap[userEntry.shortID] = userEntry.fullURL
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return shortIDMap, nil
}

func (repository *DatabaseRepository) Close() error {
	return repository.db.DBConnection.Close()
}

func (repository *DatabaseRepository) Ping(ctx context.Context) error {
	return repository.db.DBConnection.PingContext(ctx)
}

func (repository *DatabaseRepository) SaveEntries(ctx context.Context, userID int, shortIDMap ShortIDMap) error {
	tx, err := repository.db.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for shortID, fullURL := range shortIDMap {
		_, err := tx.ExecContext(ctx, insertEntrySQLBatch, shortID, fullURL, userID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
func (repository *DatabaseRepository) GetNewUserID(ctx context.Context) (userID int, err error) {
	row := repository.db.DBConnection.QueryRowContext(ctx, insertNewUser)

	err = row.Scan(&userID)

	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (repository *DatabaseRepository) DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error {
	tx, err := repository.db.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, request := range deleteRequests {
		_, err = tx.ExecContext(ctx, deleteShortLink, request.ShortIDToDelete, request.UserID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
