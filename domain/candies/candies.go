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
func (s Service) Register(ctx context.Context, candyDto domain.Candy) (*domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.Register"

	query := `INSERT INTO candies (name, price) values ($1, $2) RETURNING id, name, price`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new candy...", infra.Metadata{
		"candy": candyDto,
	})

	result, err := s.in.Db.ExecuteQuery(ctx, query, candyDto.Name, candyDto.Price)
	defer result.Close()

	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	result.Next()
	candy := domain.Candy{}
	if err := result.Scan(&candy.ID, &candy.Name, &candy.Price); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	if err := result.Err(); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &candy, nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.List"

	query := `SELECT id, name, price from candies`

	s.in.Log.Info(ctx, opName, "Listing all candies...")

	result, err := s.in.Db.ExecuteQuery(ctx, query)
	defer result.Close()

	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

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
