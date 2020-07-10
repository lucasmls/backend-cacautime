package users

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
	const opName infra.OpName = "users.NewService"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Db == nil {
		err := infra.MissingDependencyError{DependencyName: "RelationalDatabaseProvider"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// FindByEmail ...
func (s Service) FindByEmail(ctx context.Context, email string) (*domain.User, *infra.Error) {
	const opName infra.OpName = "users.Find"

	s.in.Log.Info(ctx, opName, "Fetching the user...")

	query := `
		SELECT
			u.id as id,
			u.name as name,
			u.email as email,
			u.password as password
		FROM
			users u
		WHERE email = $1
	`

	decoder := s.in.Db.Query(ctx, query, email)

	user := domain.User{}
	if err := decoder.Decode(ctx, &user); err != nil {
		return nil, errors.New(ctx, opName, err)
	}

	return &user, nil
}
