package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	connect *pgx.Conn
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	return &PostgresRepository{connect: conn}, nil
}

func (repo *PostgresRepository) Ping(ctx context.Context) error {
	return repo.connect.Ping(ctx)
}
