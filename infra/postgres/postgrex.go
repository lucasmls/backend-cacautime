package postgres

import (
	"context"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"

	// Needed only to enable postgres driver
	_ "github.com/lib/pq"
)

// ClientInput ...
type ClientInput struct {
	Log                infra.LogProvider
	ConnectionString   string
	MaxConnectionsOpen int
}

// Client ...
type Client struct {
	in ClientInput
	db *sqlx.DB
}

// NewClient ...
func NewClient(in ClientInput) (*Client, *infra.Error) {
	const opName infra.OpName = "postgres.NewClient"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.ConnectionString == "" {
		err := infra.MissingDependencyError{DependencyName: "ConnectionString"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.MaxConnectionsOpen < 1 {
		err := infra.MinimumValueError{EnvVarName: "MaxConnectionsOpen", MinimumRequired: 1}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	db, err := sqlx.Open("postgres", in.ConnectionString)
	if err != nil {
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

// Query - Executes a query that return only one row
func (c Client) Query(ctx context.Context, query string, args ...interface{}) infra.Decoder {
	const opName infra.OpName = "postgres.Query"

	c.in.Log.DebugMetadata(ctx, opName, "Executing query...", infra.Metadata{
		"query": query,
		"args":  args,
	})

	row := c.db.QueryRowxContext(ctx, query, args...)

	return decoder{row: row}
}

// QueryAll - Executes a query that returns many rows
func (c Client) QueryAll(ctx context.Context, query string, args ...interface{}) (infra.Cursor, *infra.Error) {
	const opName infra.OpName = "postgres.QueryAll"

	c.in.Log.DebugMetadata(ctx, opName, "Executing query...", infra.Metadata{
		"query": query,
		"args":  args,
	})

	rows, err := c.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.New(ctx, err, opName)
	}

	return cursor{rows: rows}, nil
}

// Execute a query that return no row(s)
func (c Client) Execute(ctx context.Context, query string, args ...interface{}) (driver.Result, *infra.Error) {
	const opName infra.OpName = "postgres.Execute"

	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.New(ctx, err, opName, infra.Metadata{
			"query": query,
			"args":  args,
		})
	}

	return result, nil
}
