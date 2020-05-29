package domain

import (
	"context"

	"github.com/lucasmls/backend-cacautime/infra"
)

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) (*Customer, *infra.Error)
	Update(context.Context, infra.ObjectID, Customer) (*Customer, *infra.Error)
	Delete(context.Context, infra.ObjectID) *infra.Error
	List(context.Context) ([]Customer, *infra.Error)
}

// DutiesRepository ...
type DutiesRepository interface {
	List(context.Context) ([]Duty, *infra.Error)
	Register(context.Context, Duty) (*Duty, *infra.Error)
	Update(context.Context, infra.ObjectID, Duty) (*Duty, *infra.Error)
	Sales(context.Context, infra.ObjectID) (*DutySales, *infra.Error)
}

// CandiesRepository ...
type CandiesRepository interface {
	Register(context.Context, Candy) (*Candy, *infra.Error)
	List(context.Context) ([]Candy, *infra.Error)
	Update(context.Context, infra.ObjectID, Candy) (*Candy, *infra.Error)
	Delete(context.Context, infra.ObjectID) *infra.Error
}

// SalesRepository ...
type SalesRepository interface {
	Register(context.Context, Sale) (*Sale, *infra.Error)
	Update(context.Context, infra.ObjectID, Sale) (*Sale, *infra.Error)
}
