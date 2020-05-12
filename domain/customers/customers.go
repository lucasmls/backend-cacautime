package customers

import (
	"context"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
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
	const opName infra.OpName = "customers.NewService"

	if in.Db == nil {
		err := infra.MissingDependencyError{DependencyName: "DatabaseClient"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Register ...
func (s Service) Register(ctx context.Context, customer domain.Customer) error {
	const opName infra.OpName = "customers.Register"

	query := `INSERT INTO customers (name, phone) values ($1, $2)`

	_, err := s.in.Db.Execute(ctx, query, customer.Name, customer.Phone)
	if err != nil {
		return errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return nil
}
