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
		`SELECT user_id, short_id, full_url, is_deleted FROM %s
		WHERE short_id = $1`, database.ShortLinksTableName)
	selectShortIDSQL = fmt.Sprintf(
		`SELECT short_id FROM %s
		WHERE full_url = $1`, database.ShortLinksTableName)
	selectUserURLs = fmt.Sprintf(
		`SELECT short_id, full_url FROM %s
		WHERE user_id = $1`, database.ShortLinksTableName)
	insertNewUserSQL = fmt.Sprintf(
		`INSERT INTO %s DEFAULT VALUES RETURNING id;`,
		database.UsersTableName)
	deleteShortLinkSQL = fmt.Sprintf(
		`UPDATE %s
		SET is_deleted = true
		WHERE short_id = $1 AND user_id = $2`,
		database.ShortLinksTableName)
)

// DatabaseRepository репозиторий, использующий БД
type DatabaseRepository struct {
	db *database.Database
}

// NewDatabaseRepository инициирует подключение к БД
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

// GetFullURL ищет в хранилище полную ссылку на ресурс по короткому ID
func (repository *DatabaseRepository) GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *ShortenedURLInfo, ok bool) {
	shortenedURLInfo = &ShortenedURLInfo{}
	row := repository.db.DBConnection.QueryRowContext(ctx, selectFullURLSQL, shortID)
	err := row.Scan(&shortenedURLInfo.UserID, &shortenedURLInfo.ShortID, &shortenedURLInfo.FullURL, &shortenedURLInfo.IsDeleted)

	if err != nil {
		return nil, false
	}

	return shortenedURLInfo, true
}

// SaveEntry сохраняет в хранилище информацию о сокращенной ссылке
func (repository *DatabaseRepository) SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error {
	stmt, err := repository.db.DBConnection.PrepareContext(ctx, insertEntrySQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, shortID, fullURL, userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
	}

	return err
}

// GetShortID возвращает shortID сокращенной ссылки
func (repository *DatabaseRepository) GetShortID(ctx context.Context, fullURL string) (shortID string, err error) {
	stmt, err := repository.db.DBConnection.PrepareContext(ctx, selectShortIDSQL)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, selectShortIDSQL, fullURL)
	err = row.Scan(&shortID)

	if err != nil {
		return "", err
	}

	return
}

// GetUserEntries возвращает сокращенный пользователем ссылки по userID
func (repository *DatabaseRepository) GetUserEntries(ctx context.Context, userID int) (shortIDMap URLMapping, err error) {
	stmt, err := repository.db.DBConnection.PrepareContext(ctx, selectUserURLs)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, selectUserURLs, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shortIDMap = make(URLMapping, 0)
	for rows.Next() {
		var userEntry struct {
			shortID string
			fullURL string
		}

		err = rows.Scan(&userEntry.shortID, &userEntry.fullURL)
		if err != nil {
			return nil, err
		}

		shortIDMap[(userEntry.shortID)] = userEntry.fullURL
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return shortIDMap, nil
}

// Close закрывает подключение к БД
func (repository *DatabaseRepository) Close() error {
	return repository.db.DBConnection.Close()
}

// Ping проверяет подключение к БД
func (repository *DatabaseRepository) Ping(ctx context.Context) error {
	return repository.db.DBConnection.PingContext(ctx)
}

// SaveEntries записывает набор сокращенных ссылок в БД
func (repository *DatabaseRepository) SaveEntries(ctx context.Context, userID int, shortIDMap URLMapping) error {
	tx, err := repository.db.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, selectUserURLs)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for shortID, fullURL := range shortIDMap {
		_, err := stmt.ExecContext(ctx, insertEntrySQLBatch, shortID, fullURL, userID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetNewUserID возвращает ID нового пользователя
func (repository *DatabaseRepository) GetNewUserID(ctx context.Context) (userID int, err error) {
	stmt, err := repository.db.DBConnection.PrepareContext(ctx, selectUserURLs)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, insertNewUserSQL)

	err = row.Scan(&userID)

	if err != nil {
		return -1, err
	}

	return userID, nil
}

// DeleteUserURLs удаляет набор сокращенных ссылок
func (repository *DatabaseRepository) DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error {
	tx, err := repository.db.DBConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, selectUserURLs)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, request := range deleteRequests {
		_, err = stmt.ExecContext(ctx, deleteShortLinkSQL, request.ShortIDToDelete, request.UserID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
