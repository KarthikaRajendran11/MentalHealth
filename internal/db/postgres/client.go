package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // dialect registers itself on init
	_ "github.com/lib/pq"
)

// Client gets or saves a profile to a database
type Client struct {
	sql *sql.DB
}

func NewClient(ctx context.Context, dbURL string) (*Client, error) {
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return &Client{
		sql: sqlDB,
	}, nil
}

func (c *Client) Write(url string) error {

	_, err := c.sql.Exec(`INSERT INTO urlhistory(url, visitTime) VALUES($1, $2)`, url, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
