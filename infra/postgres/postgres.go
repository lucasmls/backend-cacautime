package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"

	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"

	// Needed only to enable postgres driver
	_ "github.com/lib/pq"
)

// ClientInput ...
type ClientInput struct {
	Driver             string
	ConnectionString   string
	MaxConnectionsOpen int
}

// Client ...
type Client struct {
	in ClientInput
	db *sql.DB
}

// NewClient ...
func NewClient(in ClientInput) (*Client, *infra.Error) {
	// @TODO => Validar as entradas...
	const opName infra.OpName = "postgres.NewClient"

	db, err := sql.Open(in.Driver, in.ConnectionString)
	if err != nil {
		log.Panic(err)
		return nil, errors.New(err, opName, "Failed to connect into postgres.", infra.KindBadRequest)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.New(err, opName, "Failed to ping postgres.", infra.KindBadRequest)
	}

	db.SetMaxOpenConns(in.MaxConnectionsOpen)

	return &Client{
		in: in,
		db: db,
	}, nil
}

// ExecuteQuery fetch query results from the database.
func (c Client) ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// Execute a query without returning any rows.
func (c Client) Execute(ctx context.Context, query string, args ...interface{}) (driver.Result, error) {
	return c.db.Exec(query, args...)
}
