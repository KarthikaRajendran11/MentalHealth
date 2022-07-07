package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // dialect registers itself on init
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Client gets or saves a profile to a database
type Client struct {
	sql *sql.DB
}

func NewClient(ctx context.Context, dbURL string) (*Client, error) {
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open connection to DB")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping DB")
	}

	return &Client{
		sql: sqlDB,
	}, nil
}

// TBD: Do we need retry here ?
func (c *Client) Write(urlEmail ...string) error {
	_, err := c.sql.Exec(`INSERT INTO urlhistory(url, email, visitTime) VALUES($1, $2, $3)`, urlEmail[0], urlEmail[1], time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "failed to insert row")
	}
	return nil
}
