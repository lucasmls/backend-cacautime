package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// cursor ...
type cursor struct {
	rows *sqlx.Rows
}

// Next ...
func (c cursor) Next(ctx context.Context) bool {
	return c.rows.Next()
}

// Decode ...
func (c cursor) Decode(ctx context.Context, dest infra.Entity) *infra.Error {
	const opName infra.OpName = "postgres.cursor.Decode"

	if err := c.rows.StructScan(dest); err != nil {
		return errors.New(ctx, opName, err)
	}

	return nil
}

// Close ...
func (c cursor) Close(ctx context.Context) *infra.Error {
	const opName infra.OpName = "postgres.cursor.Close"

	if err := c.rows.Close(); err != nil {
		return errors.New(ctx, opName, err)
	}

	return nil
}
