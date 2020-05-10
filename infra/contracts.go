package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
)
// DatabaseClient ...
type DatabaseClient interface {
	ExecuteQuery(context.Context, string, ...interface{}) (*sql.Rows, error)
	Execute(context.Context, string, ...interface{}) (driver.Result, error)
}
