package duties

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

// Consolidate ...
func (s Service) Consolidate(ctx context.Context) (domain.ConsolidatedDuties, *infra.Error) {
	const opName infra.OpName = "duties.Consolidate"

	query := `
		SELECT
			d.id as duty_id,
			d.date as duty_date,
			d.candy_quantity as duty_candy_qtd,

			s.id as sale_id,
			s.status as sale_status,
			s.payment_method as sale_payment_method,
			
			ca.id as candy_id,
			ca.name as candy_name,
			ca.price as candy_price,
			
			cu.id as customer_id,
			cu.name as customer_name,
			cu.phone as customer_phone
		
		FROM duties d
			INNER JOIN sales s on d.id = s.duty_id
			INNER JOIN customers cu on s.customer_id = cu.id
			INNER JOIN candies ca on s.candy_id = ca.id
	`

	s.in.Log.Info(ctx, opName, "Consolidating the sales of the duties...")

	dbResult, err := s.in.Db.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer dbResult.Close()

	var sales []domain.SaleRow

	for dbResult.Next() {
		sale := domain.SaleRow{}

		if err := dbResult.Scan(
			&sale.DutyID,
			&sale.DutyDate,
			&sale.DutyQuantity,
			&sale.ID,
			&sale.Status,
			&sale.PaymentMethod,
			&sale.CandyID,
			&sale.CandyName,
			&sale.CandyPrice,
			&sale.CustomerID,
			&sale.CustomerName,
			&sale.CustomerPhone,
		); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		sales = append(sales, sale)
	}

	if err := dbResult.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	consolidatedDuties := make(domain.ConsolidatedDuties)

	for _, sale := range sales {
		var duty domain.ConsolidatedDuty

		if foundDuty, ok := consolidatedDuties[sale.DutyID]; ok {
			duty = foundDuty
		} else {
			duty = domain.ConsolidatedDuty{
				ID:       sale.DutyID,
				Date:     sale.DutyDate,
				Quantity: sale.DutyQuantity,
			}
		}

		if duty.Sales == nil {
			duty.Sales = []domain.DutySale{}
		}

		duty.Sales = append(duty.Sales, domain.DutySale{
			ID:            sale.ID,
			CandyID:       sale.CandyID,
			CandyName:     sale.CandyName,
			CandyPrice:    sale.CandyPrice,
			CustomerID:    sale.CustomerID,
			CustomerName:  sale.CustomerName,
			CustomerPhone: sale.CustomerPhone,
			PaymentMethod: sale.PaymentMethod,
			Status:        sale.Status,
		})

		duty.Subtotal += sale.CandyPrice

		if sale.Status == domain.Paid {
			duty.PaidAmount += sale.CandyPrice
		}

		if sale.Status == domain.NotPaid {
			duty.ScheduledAmount += sale.CandyPrice
		}

		consolidatedDuties[sale.DutyID] = duty
	}

	s.in.Log.InfoMetadata(ctx, opName, "Consolidated duties...", infra.Metadata{
		"duties": consolidatedDuties,
	})

	return consolidatedDuties, nil
}
