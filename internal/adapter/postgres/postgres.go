package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Host     string `envconfig:"DB_HOST"`
	Port     string `envconfig:"DB_PORT"`
	DBName   string `envconfig:"DB_NAME"`
}

type Pool struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, c Config) (*Pool, error) {
	const op = "postgres.New"

	DBURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=5s",
		c.User, c.Password, c.Host, c.Port, c.DBName)

	pool, err := pgxpool.New(ctx, DBURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to pool: %s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping pool: %s, %w", op, err)
	}

	return &Pool{pool}, nil
}

// Implement CRUD
