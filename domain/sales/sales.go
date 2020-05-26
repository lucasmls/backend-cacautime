package sales

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
	const opName infra.OpName = "sales.NewService"

	if in.Db == nil {
		err := infra.MissingDependencyError{DependencyName: "Db"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Register ...
func (s Service) Register(ctx context.Context, saleDTO domain.Sale) (*domain.Sale, *infra.Error) {
	const opName infra.OpName = "sales.Register"

	query := `INSERT INTO sales (customer_id, duty_id, candy_id, status, payment_method) values ($1, $2, $3, $4, $5) RETURNING id, customer_id as customerId, duty_id as dutyId, candy_id as candyId, status, payment_method as paymentMethod`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale": saleDTO,
	})

	decoder := s.in.Db.Query(ctx, query, saleDTO.CustomerID, saleDTO.DutyID, saleDTO.CandyID, saleDTO.Status, saleDTO.PaymentMethod)
	sale := domain.Sale{}

	if err := decoder.Decode(ctx, &sale); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &sale, nil
}
