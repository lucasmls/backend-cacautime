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

// Find ...
func (s Service) Find(ctx context.Context, customerID infra.ObjectID) (*domain.Customer, *infra.Error) {
	const opName infra.OpName = "customers.Find"

	s.in.Log.Info(ctx, opName, "Fetching the customer...")

	query := `
		SELECT
			cu.id as id,
			cu.name as name,
			cu.phone as phone
		FROM
			customers cu
		WHERE id = $1
	`

	decoder := s.in.Db.Query(ctx, query, customerID)

	customer := domain.Customer{}
	err := decoder.Decode(ctx, &customer)

	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	return &customer, nil
}

// Update ...
func (s Service) Update(ctx context.Context, customerID infra.ObjectID, customerDto domain.Customer) (*domain.Customer, *infra.Error) {
	const opName infra.OpName = "customers.Update"

	query := `UPDATE customers SET name = $1, phone = $2 WHERE id = $3 RETURNING id, name, phone`

	s.in.Log.InfoMetadata(ctx, opName, "Updating a customer...", infra.Metadata{
		"customerID": customerID,
		"dto":        customerDto,
	})

	_, err := s.Find(ctx, customerID)
	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	decoder := s.in.Db.Query(ctx, query, customerDto.Name, customerDto.Phone, customerID)

	customer := domain.Customer{}
	if err := decoder.Decode(ctx, &customer); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	return &customer, nil
}

// Delete ...
func (s Service) Delete(ctx context.Context, customerID infra.ObjectID) *infra.Error {
	const opName infra.OpName = "customers.Delete"

	query := `DELETE from customers WHERE id = $1`

	s.in.Log.InfoMetadata(ctx, opName, "Deleting a customer...", infra.Metadata{
		"customerID": customerID,
	})

	result, err := s.in.Db.Execute(ctx, query, customerID)
	if err != nil {
		return errors.New(ctx, opName, err)
	}

	affectedRowsCount, rErr := result.RowsAffected()
	if rErr != nil {
		return errors.New(ctx, opName, rErr)
	}

	if affectedRowsCount < 1 {
		return errors.New(ctx, opName, "The customer was not found.", infra.KindNotFound)
	}

	return nil
}
