package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func NewConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return pool, nil
}
