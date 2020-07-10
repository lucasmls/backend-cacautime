package auth

import (
	"context"
	"fmt"

	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ServiceInput ...
type ServiceInput struct {
	Log    infra.LogProvider
	Crypto infra.CryptoProvider
	Users  domain.UsersRepository
	JWT    infra.TokenProvider
}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) (*Service, *infra.Error) {
	const opName infra.OpName = "candies.NewService"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Crypto == nil {
		err := infra.MissingDependencyError{DependencyName: "Crypto"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Users == nil {
		err := infra.MissingDependencyError{DependencyName: "UsersRepository"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.JWT == nil {
		err := infra.MissingDependencyError{DependencyName: "TokenProvider"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Login ...
func (s Service) Login(ctx context.Context, email string, password string) (string, *infra.Error) {
	const opName infra.OpName = "auth.Login"

	user, err := s.in.Users.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New(ctx, opName, err)
	}

	if err := s.in.Crypto.Compare(ctx, user.Password, password); err != nil {
		return "", errors.New(ctx, opName, err)
	}

	jwt, err := s.in.JWT.Generate(ctx, fmt.Sprintf("%d", user.ID))
	if err != nil {
		return "", errors.New(ctx, opName, err)
	}

	return jwt, nil
}
