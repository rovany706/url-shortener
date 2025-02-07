package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Ping(ctx context.Context) error
	Close()
}

type PgDatabase struct {
	dbConnection *pgxpool.Pool
}

func InitConnection(ctx context.Context, connString string) (*PgDatabase, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	db := PgDatabase{
		dbConnection: pool,
	}
	return &db, nil
}

func (db *PgDatabase) Ping(ctx context.Context) error {
	return db.dbConnection.Ping(ctx)
}

func (db *PgDatabase) Close() {
	db.dbConnection.Close()
}
