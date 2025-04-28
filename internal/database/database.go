package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Имена таблиц
const (
	// ShortLinksTableName имя таблицы сокращенных ссылок
	ShortLinksTableName = "short_links"
	// ShortLinksTableName имя таблицы пользователей
	UsersTableName = "users"
)

var сreateTablesSQL = fmt.Sprintf(
	`DROP TABLE IF EXISTS %[2]s;
	DROP TABLE IF EXISTS %[1]s;

	CREATE TABLE %[1]s (
		id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
	);

	CREATE TABLE %[2]s (
		id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		short_id varchar(8) NOT NULL,
		full_url text UNIQUE NOT NULL,
		is_deleted boolean NOT NULL,
		user_id INT REFERENCES users(id)
	);`,
	UsersTableName, ShortLinksTableName)

// Database хранит подключение к БД
type Database struct {
	DBConnection *sql.DB
}

// InitConnection инициализирует подключение к БД
func InitConnection(ctx context.Context, connString string) (*Database, error) {
	dbConnection, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	db := Database{
		DBConnection: dbConnection,
	}
	return &db, nil
}

// EnsureCreated создает необходимые для работы таблицы
func (db *Database) EnsureCreated(ctx context.Context) error {
	result := true

	for _, tableName := range []string{UsersTableName, ShortLinksTableName} {
		ok, err := db.tableExists(ctx, tableName)
		if err != nil {
			return err
		}
		result = result && ok
	}

	if result {
		return nil
	}

	_, err := db.DBConnection.ExecContext(ctx, сreateTablesSQL)

	return err
}

func (db *Database) tableExists(ctx context.Context, tableName string) (bool, error) {
	var n int64
	err := db.DBConnection.QueryRowContext(ctx, "SELECT 1 FROM information_schema.tables where table_name = $1", tableName).Scan(&n)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
