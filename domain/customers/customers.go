package customers

import (
	"context"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ServiceInput ...
type ServiceInput struct {
	Db  infra.RelationalDatabaseProvider
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
func (s Service) Register(ctx context.Context, customerDto domain.Customer) (*domain.Customer, *infra.Error) {
	const opName infra.OpName = "customers.Register"

	query := `INSERT INTO customers (name, phone) values ($1, $2) RETURNING id, name, phone`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new customer...", infra.Metadata{
		"customer": customerDto,
	})

	decoder := s.in.Db.Query(ctx, query, customerDto.Name, customerDto.Phone)

	customer := domain.Customer{}
	if err := decoder.Decode(ctx, &customer); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	return &customer, nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Customer, *infra.Error) {
	const opName infra.OpName = "customers.List"

	query := `SELECT id, name, phone from customers`

	s.in.Log.Info(ctx, opName, "Listing all customers...")

	cursor, err := s.in.Db.QueryAll(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer cursor.Close(ctx)

	var customers []domain.Customer

	for cursor.Next(ctx) {
		customer := domain.Customer{}
		if err := cursor.Decode(ctx, &customer); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
