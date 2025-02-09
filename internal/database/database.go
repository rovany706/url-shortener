package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const TableName = "short_links"

var CreateTableSql = fmt.Sprintf(
	`CREATE TABLE %s (
	id SERIAL PRIMARY KEY,
	short_id varchar(8),
	full_url text)`,
	TableName)

type Database interface {
	EnsureCreated(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PingContext(ctx context.Context) error
	Close()
}

type SqlDatabase struct {
	dbConnection *sql.DB
}

func InitConnection(ctx context.Context, connString string) (*SqlDatabase, error) {
	dbConnection, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	db := SqlDatabase{
		dbConnection: dbConnection,
	}
	return &db, nil
}

func (db *SqlDatabase) EnsureCreated(ctx context.Context) error {
	ok, err := db.tableExists(ctx, "short_links")
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	_, err = db.ExecContext(ctx, CreateTableSql)

	return err
}

func (db *SqlDatabase) tableExists(ctx context.Context, tableName string) (bool, error) {
	var n int64
	err := db.dbConnection.QueryRowContext(ctx, "SELECT 1 FROM information_schema.tables where table_name = $1", tableName).Scan(&n)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *SqlDatabase) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.dbConnection.ExecContext(ctx, query, args...)
}

func (db *SqlDatabase) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.dbConnection.QueryContext(ctx, query, args...)
}
func (db *SqlDatabase) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.dbConnection.QueryRowContext(ctx, query, args...)
}

func (db *SqlDatabase) PingContext(ctx context.Context) error {
	return db.dbConnection.PingContext(ctx)
}

func (db *SqlDatabase) Close() {
	db.dbConnection.Close()
}
