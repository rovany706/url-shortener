package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const TableName = "short_links"

var CreateTableSQL = fmt.Sprintf(
	`CREATE TABLE %s (
	id SERIAL PRIMARY KEY,
	short_id varchar(8),
	full_url text UNIQUE)`,
	TableName)

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
	ok, err := db.tableExists(ctx, TableName)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	_, err = db.DBConnection.ExecContext(ctx, CreateTableSQL)

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
