package sales

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
	const opName infra.OpName = "sales.NewService"

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
func (s Service) Register(ctx context.Context, sale domain.Sale) *infra.Error {
	const opName infra.OpName = "sales.Register"

	query := `INSERT INTO sales (customer_id, duty_id, candy_id, status, payment_method) values ($1, $2, $3, $4, $5)`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale": sale,
	})

	_, err := s.in.Db.Execute(ctx, query, sale.CustomerID, sale.DutyID, sale.CandyID, sale.Status, sale.PaymentMethod)
	if err != nil {
		return errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return nil
}
