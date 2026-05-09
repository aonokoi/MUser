package postgres

import (
	"context"
	"fmt"

	"goproj/MUser/internal/domain"

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

func (p *Pool) CreateUser(ctx context.Context, name domain.Name, email domain.Email) (int, error) {
	const op = "postgres.CreateUser"

	sql := `
	INSERT INTO users(name, email)
	VALUES($1, $2)
	RETURNING id
	`

	var id int

	err := p.pool.QueryRow(ctx, sql, name, email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable to create user: %s: %w", op, err)
	}

	return id, nil
}

func (p *Pool) ReadUser(ctx context.Context, id int) (domain.User, error) {
	const op = "postgres.ReadUser"

	sql := `SELECT * FROM users WHERE id = $1`

	var user domain.User

	err := p.pool.QueryRow(ctx, sql, id).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
	if err != nil {
		return user, fmt.Errorf("unable to read info about user: %s: %w", op, err)
	}

	return user, nil
}
