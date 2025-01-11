package di

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

type postgreSQL struct {
	client *sql.DB
}

type postgreSQLOptions struct {
	username     string
	password     string
	host         string
	database     string
	timeout      string
	maxIdleConns int
	maxOpenConns int
	maxLifetime  time.Duration
}

func newPostgreSQL(opts postgreSQLOptions) (*postgreSQL, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?statement_timeout=%s&sslmode=disable",
		opts.username,
		opts.password,
		opts.host,
		opts.database,
		opts.timeout,
	)

	log.Printf("[di.newPostgreSQL] connecting to PostgreSQL database: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opts.maxIdleConns)
	db.SetMaxOpenConns(opts.maxOpenConns)
	db.SetConnMaxLifetime(opts.maxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	slog.Info("[di.newPostgreSQL] PostgreSQL database connected",
		slog.String("database", opts.database),
	)

	return &postgreSQL{db}, nil
}
