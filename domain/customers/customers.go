package customers

import (
	"context"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
)

// ServiceInput ...
type ServiceInput struct {
	Db infra.DatabaseClient
}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) (*Service, *infra.Error) {
	// @TODO => Validar as entradas...
	const opName infra.OpName = "customers.NewService"

	return &Service{
		in: in,
	}, nil
}

// Register ...
func (s Service) Register(ctx context.Context, customer domain.Customer) error {
	query := `INSERT INTO customers (name, phone) values ($1, $2)`

	_, err := s.in.Db.Execute(ctx, query, customer.Name, customer.Phone)
	if err != nil {
		return err
	}

	return nil
}
