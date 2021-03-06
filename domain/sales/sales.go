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

	query := "INSERT INTO sales (customer_id, candy_id, status, payment_method, date) values ($1, $2, $3, $4, $5) RETURNING id, customer_id as customerId, candy_id as candyId, status, payment_method as paymentMethod, date as date"

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale": saleDTO,
	})

	s.in.Log.DebugMetadata(ctx, opName, "Registering a new sale...", infra.Metadata{
		"sale":  saleDTO,
		"query": query,
	})

	decoder := s.in.Db.Query(ctx, query, saleDTO.CustomerID, saleDTO.CandyID, saleDTO.Status, saleDTO.PaymentMethod, saleDTO.Date)
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
			sa.status as status,
			sa.date::text as date
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
			date::text,
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

// Months ...
func (s Service) Months(ctx context.Context) ([]domain.Month, *infra.Error) {
	const opName infra.OpName = "sales.Months"

	query := `
		WITH months_with_sales AS (
			SELECT 
				trim(to_char(date, 'Month')) as month,
				trim(to_char(date, 'MM')) as number,
				trim(to_char(date, 'YYYY')) as year
			FROM sales
			GROUP BY 1, 2, 3
		)
		SELECT *
		FROM months_with_sales
		ORDER BY date(concat('01/', number, '/', year)) DESC;`

	s.in.Log.Info(ctx, opName, "Listing months that has sales...")

	cursor, err := s.in.Db.QueryAll(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer cursor.Close(ctx)

	months := []domain.Month{}

	for cursor.Next(ctx) {
		group := domain.Month{}
		if err := cursor.Decode(ctx, &group); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		months = append(months, group)
	}

	return months, nil
}

// MonthSales ...
func (s Service) MonthSales(ctx context.Context, month int, year int) (*domain.MonthSales, *infra.Error) {
	const opName infra.OpName = "sales.MonthSales"

	s.in.Log.Info(ctx, opName, "Fetching the month sales")

	query := `
		SELECT
			s.id as id,
			s.status as status,
			s.payment_method as paymentMethod,
			trim(to_char(s.date, 'DD/MM/YYYY')) as date,
			
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
			EXTRACT(MONTH FROM s.date) = $1 and EXTRACT(YEAR FROM s.date) = $2
		ORDER BY s.created_at;
	`

	cursor, dbErr := s.in.Db.QueryAll(ctx, query, month, year)
	if dbErr != nil {
		return nil, errors.New(ctx, opName, dbErr, infra.KindBadRequest)
	}

	defer cursor.Close(ctx)

	monthSales := domain.MonthSales{
		Subtotal:        0,
		PaidAmount:      0,
		ScheduledAmount: 0,
		Sales:           []domain.MonthSale{},
	}

	for cursor.Next(ctx) {
		sale := domain.MonthSale{}
		if err := cursor.Decode(ctx, &sale); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
		}

		monthSales.Sales = append(monthSales.Sales, sale)
		monthSales.Subtotal += sale.CandyPrice

		if sale.Status == domain.Paid {
			monthSales.PaidAmount += sale.CandyPrice
		}

		if sale.Status == domain.NotPaid {
			monthSales.ScheduledAmount += sale.CandyPrice
		}
	}

	return &monthSales, nil
}
