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

// Sales ...
func (s Service) Sales(ctx context.Context) (domain.DutiesResult, *infra.Error) {
	const opName infra.OpName = "duties.Sales"

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

	s.in.Log.Info(ctx, opName, "Listing the sales of the duties...")

	result, err := s.in.Db.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer result.Close()

	var dutiesSalesResult []domain.RawDutyResult

	for result.Next() {
		duty := domain.RawDutyResult{}
		if err := result.Scan(
			&duty.ID,
			&duty.Date,
			&duty.Quantity,
			&duty.SaleID,
			&duty.SaleStatus,
			&duty.SalePaymentMethod,
			&duty.CandyID,
			&duty.CandyName,
			&duty.CandyPrice,
			&duty.CustomerID,
			&duty.CustomerName,
			&duty.CustomerPhone,
		); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		dutiesSalesResult = append(dutiesSalesResult, duty)
	}

	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	dutiesResult := make(domain.DutiesResult)

	for _, duty := range dutiesSalesResult {
		var dutyResult domain.DutyResult

		if foundResult, ok := dutiesResult[duty.ID]; !ok {
			dutyResult = domain.DutyResult{
				ID:       duty.ID,
				Date:     duty.Date,
				Quantity: duty.Quantity,
			}
		} else {
			dutyResult = foundResult
		}

		if dutyResult.Sales == nil {
			dutyResult.Sales = []domain.SaleResult{}
		}

		dutyResult.Sales = append(dutyResult.Sales, domain.SaleResult{
			ID:            duty.SaleID,
			CandyID:       duty.CandyID,
			CandyName:     duty.CandyName,
			CandyPrice:    duty.CandyPrice,
			CustomerID:    duty.CustomerID,
			CustomerName:  duty.CustomerName,
			CustomerPhone: duty.CustomerPhone,
			PaymentMethod: duty.SalePaymentMethod,
			Status:        duty.SaleStatus,
		})

		dutyResult.Subtotal = 0
		dutyResult.PaidAmount = 0
		dutyResult.ScheduledAmount = 0

		dutiesResult[duty.ID] = dutyResult
	}

	s.in.Log.InfoMetadata(ctx, opName, "Resultado do plantão...", infra.Metadata{
		"result": dutiesResult,
	})

	return dutiesResult, nil
}
