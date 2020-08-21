package domain

import (
	"context"

	"github.com/lucasmls/backend-cacautime/infra"
)

// UsersRepository ...
type UsersRepository interface {
	FindByEmail(context.Context, string) (*User, *infra.Error)
}

// AuthRepository ...
type AuthRepository interface {
	Login(context.Context, string, string) (string, *infra.Error)
}

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) (*Customer, *infra.Error)
	Update(context.Context, infra.ObjectID, Customer) (*Customer, *infra.Error)
	Delete(context.Context, infra.ObjectID) *infra.Error
	List(context.Context) ([]Customer, *infra.Error)
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
	Delete(context.Context, infra.ObjectID) *infra.Error
	Months(context.Context) ([]Month, *infra.Error)
	MonthSales(context.Context, int, int) (*MonthSales, *infra.Error)
}
