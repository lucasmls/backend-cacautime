package duties

import (
	"context"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ServiceInput ...
type ServiceInput struct {
	Log infra.LogProvider
	Db  infra.RelationalDatabaseProvider
}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) (*Service, *infra.Error) {
	const opName infra.OpName = "duties.NewService"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Db == nil {
		err := infra.MissingDependencyError{DependencyName: "Db"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Register ...
func (s Service) Register(ctx context.Context, dutyDTO domain.Duty) (*domain.Duty, *infra.Error) {
	const opName infra.OpName = "duties.Register"

	query := `INSERT INTO duties (date, candy_quantity) values ($1, $2) RETURNING id, date, candy_quantity as candyQuantity`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new duty...", infra.Metadata{
		"duty": dutyDTO,
	})

	decoder := s.in.Db.Query(ctx, query, dutyDTO.Date, dutyDTO.CandyQuantity)

	duty := domain.Duty{}

	if err := decoder.Decode(ctx, &duty); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &duty, nil
}

// Update ...
func (s Service) Update(ctx context.Context, dutyID infra.ObjectID, dutyDTO domain.Duty) (*domain.Duty, *infra.Error) {
	const opName infra.OpName = "duties.Update"

	query := `UPDATE duties SET date = $1, candy_quantity = $2 WHERE id = $3 RETURNING id, date, candy_quantity as candyQuantity`

	s.in.Log.InfoMetadata(ctx, opName, "Updating a duty...", infra.Metadata{
		"dutyID": dutyID,
		"dto":    dutyDTO,
	})

	_, err := s.Find(ctx, dutyID)
	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	decoder := s.in.Db.Query(ctx, query, dutyDTO.Date, dutyDTO.CandyQuantity, dutyID)

	duty := domain.Duty{}
	if err := decoder.Decode(ctx, &duty); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	return &duty, nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Duty, *infra.Error) {
	const opName infra.OpName = "duties.List"

	query := `SELECT id, date, candy_quantity as candyQuantity from duties`

	s.in.Log.Info(ctx, opName, "Listing all duties...")

	cursor, err := s.in.Db.QueryAll(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer cursor.Close(ctx)

	var duties []domain.Duty

	for cursor.Next(ctx) {
		duty := domain.Duty{}
		if err := cursor.Decode(ctx, &duty); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		duties = append(duties, duty)
	}

	return duties, nil
}

// Find ...
func (s Service) Find(ctx context.Context, dutyID infra.ObjectID) (*domain.Duty, *infra.Error) {
	const opName infra.OpName = "duties.Find"

	query := `
		SELECT
			du.id as id,
			du.date as date,
			du.candy_quantity as candyQuantity
		FROM
			duties du
		WHERE id = $1
	`

	s.in.Log.Info(ctx, opName, "Fetching the duty...")

	decoder := s.in.Db.Query(ctx, query, dutyID)

	duty := domain.Duty{}
	err := decoder.Decode(ctx, &duty)

	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	return &duty, nil
}

// Sales ...
func (s Service) Sales(ctx context.Context, dutyID infra.ObjectID) (*domain.DutySales, *infra.Error) {
	const opName infra.OpName = "duties.Sales"

	s.in.Log.Info(ctx, opName, "Fetching the duty sales...")

	duty, err := s.Find(ctx, dutyID)
	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	query := `
		SELECT
			s.id as id,
			s.status as status,
			s.payment_method as paymentMethod,
			
			cu.id as customerId,
			cu.name as customerName,

			ca.id as candyId,
			ca.name as candyName,
			ca.price as candyPrice
		FROM
			sales s
			INNER JOIN customers cu ON s.customer_id = cu.id
			INNER JOIN candies ca ON s.candy_id = ca.id
		WHERE
			s.duty_id = $1
	`

	cursor, dbErr := s.in.Db.QueryAll(ctx, query, dutyID)
	if dbErr != nil {
		return nil, errors.New(ctx, opName, dbErr, infra.KindBadRequest)
	}

	defer cursor.Close(ctx)

	dutySales := domain.DutySales{
		ID:              duty.ID,
		Date:            duty.Date,
		Quantity:        duty.CandyQuantity,
		PaidAmount:      0,
		ScheduledAmount: 0,
		Subtotal:        0,
		Sales:           []domain.DutySale{},
	}

	for cursor.Next(ctx) {
		sale := domain.DutySale{}
		if err := cursor.Decode(ctx, &sale); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
		}

		dutySales.Sales = append(dutySales.Sales, sale)

		dutySales.Subtotal += sale.CandyPrice

		if sale.Status == domain.Paid {
			dutySales.PaidAmount += sale.CandyPrice
		}

		if sale.Status == domain.NotPaid {
			dutySales.ScheduledAmount += sale.CandyPrice
		}
	}

	return &dutySales, nil
}

// Delete ...
func (s Service) Delete(ctx context.Context, dutyID infra.ObjectID) *infra.Error {
	const opName infra.OpName = "duties.Delete"

	query := `DELETE from duties WHERE id = $1`

	s.in.Log.InfoMetadata(ctx, opName, "Deleting a duty...", infra.Metadata{
		"dutyID": dutyID,
	})

	result, err := s.in.Db.Execute(ctx, query, dutyID)
	if err != nil {
		return errors.New(ctx, opName, err)
	}

	affectedRowsCount, rErr := result.RowsAffected()
	if rErr != nil {
		return errors.New(ctx, opName, rErr)
	}

	if affectedRowsCount < 1 {
		return errors.New(ctx, opName, "The duty was not found.", infra.KindNotFound)
	}

	return nil
}
