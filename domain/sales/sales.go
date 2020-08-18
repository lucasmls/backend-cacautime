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

	query := "INSERT INTO sales (customer_id, candy_id, status, payment_method) values ($1, $2, $3, $4) RETURNING id, customer_id as customerId, candy_id as candyId, status, payment_method as paymentMethod"

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale": saleDTO,
	})

	decoder := s.in.Db.Query(ctx, query, saleDTO.CustomerID, saleDTO.CandyID, saleDTO.Status, saleDTO.PaymentMethod)
	sale := domain.Sale{}

	if err := decoder.Decode(ctx, &sale); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &sale, nil
}

// Find ...
func (s Service) Find(ctx context.Context, saleID infra.ObjectID) (*domain.Sale, *infra.Error) {
	const opName infra.OpName = "sales.Find"

	s.in.Log.Info(ctx, opName, "Fetching the sale...")

	query := `
		SELECT
			sa.id as id,
			sa.customer_id as customerId,
			sa.candy_id as candyId,
			sa.payment_method as paymentMethod,
			sa.status as status
		FROM
			sales sa
		WHERE id = $1
	`

	decoder := s.in.Db.Query(ctx, query, saleID)

	sale := domain.Sale{}
	err := decoder.Decode(ctx, &sale)

	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	return &sale, nil
}

// Update ...
func (s Service) Update(ctx context.Context, saleID infra.ObjectID, saleDTO domain.Sale) (*domain.Sale, *infra.Error) {
	const opName infra.OpName = "sales.Update"

	query := `
		UPDATE sales SET
			status = $1,
			payment_method = $2
		WHERE id = $3 RETURNING
			id,
			customer_id as customerId,
			candy_id as candyId,
			status,
			payment_method as paymentMethod
	`

	s.in.Log.InfoMetadata(ctx, opName, "Updating a sale...", infra.Metadata{
		"saleID": saleID,
		"dto":    saleDTO,
	})

	_, err := s.Find(ctx, saleID)
	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	decoder := s.in.Db.Query(ctx, query, saleDTO.Status, saleDTO.PaymentMethod, saleID)

	sale := domain.Sale{}
	if err := decoder.Decode(ctx, &sale); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	return &sale, nil
}

// Delete ...
func (s Service) Delete(ctx context.Context, saleID infra.ObjectID) *infra.Error {
	const opName infra.OpName = "sales.Delete"

	query := `DELETE from sales WHERE id = $1`

	s.in.Log.InfoMetadata(ctx, opName, "Deleting a sale...", infra.Metadata{
		"saleID": saleID,
	})

	result, err := s.in.Db.Execute(ctx, query, saleID)
	if err != nil {
		return errors.New(ctx, opName, err)
	}

	affectedRowsCount, rErr := result.RowsAffected()
	if rErr != nil {
		return errors.New(ctx, opName, rErr)
	}

	if affectedRowsCount < 1 {
		return errors.New(ctx, opName, "The sale was not found.", infra.KindNotFound)
	}

	return nil
}
