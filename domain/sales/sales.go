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
func (s Service) Register(ctx context.Context, saleDTO domain.Sale) (*domain.Sale, *infra.Error) {
	const opName infra.OpName = "sales.Register"

	query := `INSERT INTO sales (customer_id, duty_id, candy_id, status, payment_method) values ($1, $2, $3, $4, $5) RETURNING id, customer_id, duty_id, candy_id, status, payment_method`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale": saleDTO,
	})

	result, err := s.in.Db.ExecuteQuery(ctx, query, saleDTO.CustomerID, saleDTO.DutyID, saleDTO.CandyID, saleDTO.Status, saleDTO.PaymentMethod)
	defer result.Close()

	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	result.Next()
	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	sale := domain.Sale{}
	if err := result.Scan(&sale.ID, &sale.CustomerID, &sale.DutyID, &sale.CandyID, &sale.Status, &sale.PaymentMethod); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &sale, nil
}
