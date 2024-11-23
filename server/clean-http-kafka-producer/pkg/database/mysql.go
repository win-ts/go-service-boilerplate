// Package database provides the database connection and operations
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQL represents the MySQL database client
type MySQL struct {
	Client *sql.DB
}

// MySQLOptions represents the options for the MySQL database
type MySQLOptions struct {
	Username     string
	Password     string
	Host         string
	Database     string
	Timeout      time.Duration
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// NewMySQL creates a new MySQL database client
func NewMySQL(opts MySQLOptions) (*MySQL, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?checkConnLiveness=false&loc=Local&parseTime=true&readTimeout=%s&timeout=%s&writeTimeout=%s&maxAllowedPacket=0",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		opts.Timeout,
		opts.Timeout,
		opts.Timeout,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetConnMaxLifetime(opts.MaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	slog.Info("[pkg.NewMySQL] connecting to MySQL database",
		slog.String("database", opts.Database),
	)

	return &MySQL{db}, nil
}

// Close closes the MySQL database connection
func (db *MySQL) Close() error {
	return db.Client.Close()
}

// Ping pings the MySQL database
func (db *MySQL) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Client.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
