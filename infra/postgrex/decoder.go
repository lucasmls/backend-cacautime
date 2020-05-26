package postgrex

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

type decoder struct {
	row *sqlx.Row
}

// Decode ...
func (d decoder) Decode(ctx context.Context, dest infra.Entity) *infra.Error {
	const opName infra.OpName = "postgres.decoder.Decode"

	err := d.row.StructScan(dest)
	if err != nil && err == sql.ErrNoRows {
		return errors.New(ctx, opName, err, infra.KindNotFound)
	}

	if err != nil {
		return errors.New(ctx, opName, err)
	}

	return nil
}
