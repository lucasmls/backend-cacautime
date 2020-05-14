package domain

import (
	"context"

	"github.com/lucasmls/backend-cacautime/infra"
)

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) *infra.Error
}

// DutiesRepository ...
type DutiesRepository interface {
	List(context.Context) ([]Duty, *infra.Error)
	Register(context.Context, Duty) *infra.Error
	Sales(context.Context) (DutiesResult, *infra.Error)
}

// CandiesRepository ...
type CandiesRepository interface {
	Register(context.Context, Candy) *infra.Error
}
