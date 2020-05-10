package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
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
func NewClient(in ClientInput) (*Client, error) {
	db, err := sql.Open(in.Driver, in.ConnectionString)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
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
	fmt.Println("======")
	fmt.Println(args...)
	return c.db.Exec(query, args...)
}
