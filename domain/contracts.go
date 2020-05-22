package domain

import (
	"context"

	"github.com/lucasmls/backend-cacautime/infra"
)

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) (*Customer, *infra.Error)
	List(context.Context) ([]Customer, *infra.Error)
}

// DutiesRepository ...
type DutiesRepository interface {
	List(context.Context) ([]Duty, *infra.Error)
	Register(context.Context, Duty) *infra.Error
	Consolidate(context.Context) (ConsolidatedDuties, *infra.Error)
}

// CandiesRepository ...
type CandiesRepository interface {
	Register(context.Context, Candy) (*Candy, *infra.Error)
	List(context.Context) ([]Candy, *infra.Error)
}

// SalesRepository ...
type SalesRepository interface {
	Register(context.Context, Sale) *infra.Error
}
