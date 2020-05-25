package duties

import (
	"context"
	"database/sql"
	"fmt"

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
	const opName infra.OpName = "duties.NewService"

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
func (s Service) Register(ctx context.Context, duty domain.Duty) *infra.Error {
	const opName infra.OpName = "duties.Register"

	query := `INSERT INTO duties (date, candy_quantity) values ($1, $2)`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new duty...", infra.Metadata{
		"duty": duty,
	})

	_, err := s.in.Db.Execute(ctx, query, duty.Date, duty.CandyQuantity)
	if err != nil {
		return errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Duty, *infra.Error) {
	const opName infra.OpName = "duties.List"

	query := `SELECT id, date, candy_quantity from duties`

	s.in.Log.Info(ctx, opName, "Listing all duties...")

	result, err := s.in.Db.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer result.Close()

	var duties []domain.Duty

	for result.Next() {
		duty := domain.Duty{}
		if err := result.Scan(&duty.ID, &duty.Date, &duty.CandyQuantity); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		duties = append(duties, duty)
	}

	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
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
			du.candy_quantity as quantity
		FROM
			duties du
		WHERE id = $1
	`

	s.in.Log.Info(ctx, opName, "Fetching the duty...")

	duty := domain.Duty{}
	err := s.in.Db.ExecuteQueryRow(ctx, query, dutyID).Scan(
		&duty.ID,
		&duty.Date,
		&duty.CandyQuantity,
	)

	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New(ctx, opName, err, infra.KindNotFound)
	}

	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	fmt.Println("===")
	fmt.Println(duty)
	fmt.Println("===")

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
			s.payment_method as payment_method,
			
			cu.id as customer_id,
			cu.name as customer_name,

			ca.id as candy_id,
			ca.name as candy_name,
			ca.price as candy_price
		FROM
			sales s
			INNER JOIN customers cu ON s.customer_id = cu.id
			INNER JOIN candies ca ON s.candy_id = ca.id
		WHERE
			s.duty_id = $1
	`

	dbResult, dbErr := s.in.Db.ExecuteQuery(ctx, query, dutyID)
	if dbErr != nil {
		return nil, errors.New(ctx, opName, dbErr, infra.KindBadRequest)
	}

	defer dbResult.Close()

	dutySales := domain.DutySales{
		ID:              duty.ID,
		Date:            duty.Date,
		Quantity:        duty.CandyQuantity,
		PaidAmount:      0,
		ScheduledAmount: 0,
		Subtotal:        0,
		Sales:           []domain.DutySale{},
	}

	for dbResult.Next() {
		sale := domain.DutySale{}

		if err := dbResult.Scan(
			&sale.ID,
			&sale.Status,
			&sale.PaymentMethod,
			&sale.CustomerID,
			&sale.CustomerName,
			&sale.CandyID,
			&sale.CandyName,
			&sale.CandyPrice,
		); err != nil {
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
