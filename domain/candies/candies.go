package candies

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
	const opName infra.OpName = "candies.NewService"

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
func (s Service) Register(ctx context.Context, candy domain.Candy) *infra.Error {
	const opName infra.OpName = "candies.Register"

	query := `INSERT INTO candies (name, price) values ($1, $2)`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new candy...", infra.Metadata{
		"candy": candy,
	})

	_, err := s.in.Db.Execute(ctx, query, candy.Name, candy.Price)
	if err != nil {
		return errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.List"

	query := `SELECT id, name, price from candies`

	s.in.Log.Info(ctx, opName, "Listing all candies...")

	result, err := s.in.Db.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	defer result.Close()

	var candies []domain.Candy

	for result.Next() {
		candy := domain.Candy{}
		if err := result.Scan(&candy.ID, &candy.Name, &candy.Price); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		candies = append(candies, candy)
	}

	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return candies, nil
}
