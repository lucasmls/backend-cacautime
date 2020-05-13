package domain

import "context"

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) error
}

// DutiesRepository ...
type DutiesRepository interface {
	Register(context.Context, Duty) error
}
