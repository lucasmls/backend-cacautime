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
func (s Service) Register(ctx context.Context, duty domain.Duty) error {
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
func (s Service) List(ctx context.Context) ([]domain.Duty, error) {
	const opName infra.OpName = "duties.List"

	query := `SELECT * from duties`

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
