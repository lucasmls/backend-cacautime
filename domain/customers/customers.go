package customers

import (
	"context"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ServiceInput ...
type ServiceInput struct {
	Db  infra.DatabaseClient
	Log infra.LogProvider
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

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "LogProvider"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Register ...
func (s Service) Register(ctx context.Context, customer domain.Customer) *infra.Error {
	const opName infra.OpName = "customers.Register"

	query := `INSERT INTO customers (name, phone) values ($1, $2)`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new customer...", infra.Metadata{
		"customer": customer,
	})

	_, err := s.in.Db.Execute(ctx, query, customer.Name, customer.Phone)
	if err != nil {
		return errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Customer, *infra.Error) {
	const opName infra.OpName = "customers.List"

	query := `SELECT id, name, phone from customers`

	s.in.Log.Info(ctx, opName, "Listing all customers...")

	result, err := s.in.Db.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer result.Close()

	var customers []domain.Customer

	for result.Next() {
		customer := domain.Customer{}
		if err := result.Scan(&customer.ID, &customer.Name, &customer.Phone); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		customers = append(customers, customer)
	}

	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return customers, nil
}
