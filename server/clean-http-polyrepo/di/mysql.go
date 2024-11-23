package di

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type mySQL struct {
	client *sql.DB
}

type mySQLOptions struct {
	username     string
	password     string
	host         string
	database     string
	timeout      time.Duration
	maxIdleConns int
	maxOpenConns int
	maxLifetime  time.Duration
}

func newMySQL(opts mySQLOptions) (*mySQL, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?checkConnLiveness=false&loc=Local&parseTime=true&readTimeout=%s&timeout=%s&writeTimeout=%s&maxAllowedPacket=0",
		opts.username,
		opts.password,
		opts.host,
		opts.database,
		opts.timeout,
		opts.timeout,
		opts.timeout,
	)

	db, err := sql.Open("mysql", dsn)
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

	slog.Info("[di.newMySQL] connecting to MySQL database",
		slog.String("database", opts.database),
	)

	return &mySQL{db}, nil
}
