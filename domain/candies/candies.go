package candies

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
	const opName infra.OpName = "candies.NewService"

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
func (s Service) Register(ctx context.Context, candyDto domain.Candy) (*domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.Register"

	query := `INSERT INTO candies (name, price) values ($1, $2) RETURNING id, name, price`

	s.in.Log.InfoMetadata(ctx, opName, "Registering a new candy...", infra.Metadata{
		"candy": candyDto,
	})

	decoder := s.in.Db.Query(ctx, query, candyDto.Name, candyDto.Price)
	candy := domain.Candy{}
	if err := decoder.Decode(ctx, &candy); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	return &candy, nil
}

// List ...
func (s Service) List(ctx context.Context) ([]domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.List"

	query := `SELECT id, name, price from candies`

	s.in.Log.Info(ctx, opName, "Listing all candies...")

	cursor, err := s.in.Db.QueryAll(ctx, query)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
	}

	var candies []domain.Candy

	for cursor.Next(ctx) {
		candy := domain.Candy{}
		if err := cursor.Decode(ctx, &candy); err != nil {
			return nil, errors.New(ctx, opName, err, infra.KindUnexpected)
		}

		candies = append(candies, candy)
	}

	return candies, nil
}

// Find ...
func (s Service) Find(ctx context.Context, candyID infra.ObjectID) (*domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.Find"

	s.in.Log.Info(ctx, opName, "Fetching the candy...")

	query := `
		SELECT
			ca.id as id,
			ca.name as name,
			ca.price as price
		FROM
			candies ca
		WHERE id = $1
	`

	decoder := s.in.Db.Query(ctx, query, candyID)

	candy := domain.Candy{}
	err := decoder.Decode(ctx, &candy)

	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	return &candy, nil
}

// Update ...
func (s Service) Update(ctx context.Context, candyID infra.ObjectID, candyDTO domain.Candy) (*domain.Candy, *infra.Error) {
	const opName infra.OpName = "candies.Update"

	query := `UPDATE candies SET name = $1, price = $2 WHERE id = $3 RETURNING id, name, price`

	s.in.Log.InfoMetadata(ctx, opName, "Updating a customer...", infra.Metadata{
		"dutyID": candyID,
		"dto":    candyDTO,
	})

	_, err := s.Find(ctx, candyID)
	if err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	decoder := s.in.Db.Query(ctx, query, candyDTO.Name, candyDTO.Price, candyID)

	candy := domain.Candy{}
	if err := decoder.Decode(ctx, &candy); err != nil {
		return nil, errors.New(ctx, opName, err, infra.KindBadRequest)
	}

	return &candy, nil
}
