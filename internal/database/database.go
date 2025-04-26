package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	ShortLinksTableName = "short_links"
	UsersTableName      = "users"
)

var CreateTablesSQL = fmt.Sprintf(
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

type Database struct {
	DBConnection *sql.DB
}

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

	_, err := db.DBConnection.ExecContext(ctx, CreateTablesSQL)

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
