package domain

import "context"

// CustomersRepository ...
type CustomersRepository interface {
	Register(context.Context, Customer) error
}
