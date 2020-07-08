package bcrypt

import (
	"context"

	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
	"golang.org/x/crypto/bcrypt"
)

// ClientInput ...
type ClientInput struct {
	Log infra.LogProvider
}

// Client ...
type Client struct {
	in   ClientInput
	Cost int
}

// NewClient ...
func NewClient(in ClientInput) (*Client, *infra.Error) {
	const opName infra.OpName = "bcrypt.NewClient"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Client{
		in:   in,
		Cost: bcrypt.DefaultCost,
	}, nil
}

// Hash ...
func (c Client) Hash(ctx context.Context, value string) ([]byte, *infra.Error) {
	const opName infra.OpName = "bcrypt.Hash"

	c.in.Log.Debug(ctx, opName, "Hashing...")

	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), c.Cost)
	if err != nil {
		return nil, errors.New(ctx, err, opName)
	}

	return hashedValue, nil
}

// Compare ...
func (c Client) Compare(ctx context.Context, value string, valueToCompare string) *infra.Error {
	const opName infra.OpName = "bcrypt.Compare"

	c.in.Log.Debug(ctx, opName, "Comparing bcrypted values...")

	if err := bcrypt.CompareHashAndPassword([]byte(value), []byte(valueToCompare)); err != nil {
		return errors.New(ctx, err, opName)
	}

	return nil
}
