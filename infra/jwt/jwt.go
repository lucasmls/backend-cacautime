package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ClientInput ...
type ClientInput struct {
	Log    infra.LogProvider
	Secret string
	TTL    int
}

// Client ...
type Client struct {
	in ClientInput
}

// NewClient ...
func NewClient(in ClientInput) (*Client, *infra.Error) {
	const opName infra.OpName = "jwt.NewClient"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.Secret == "" {
		err := infra.MissingDependencyError{DependencyName: "Secret"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.TTL < 1 {
		err := infra.MinimumValueError{EnvVarName: "TTL", MinimumRequired: 1}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Client{
		in: in,
	}, nil
}

// Generate ...
func (c Client) Generate(ctx context.Context, userID string) (string, *infra.Error) {
	const opName infra.OpName = "jwt.Generate"

	jwtInstance := jwt.New(jwt.SigningMethodHS256)
	claims := jwtInstance.Claims.(jwt.MapClaims)

	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(c.in.TTL)).Unix()

	token, err := jwtInstance.SignedString([]byte(c.in.Secret))
	if err != nil {
		return "", errors.New(ctx, err, opName)
	}

	return token, nil
}

// Validate ...
func (c Client) Validate(ctx context.Context, tokenToValidate string) (*infra.DecodedJWT, *infra.Error) {
	const opName infra.OpName = "jwt.Validate"

	token, err := jwt.Parse(tokenToValidate, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ctx, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]), opName)
		}

		return []byte(c.in.Secret), nil
	})

	if err != nil {
		return nil, errors.New(ctx, err, opName)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, errors.New(ctx, "Invalid token", opName)
	}

	decodedJWT := infra.DecodedJWT{
		UserID: claims["userID"].(string),
		Exp:    int64(claims["exp"].(float64)),
	}

	return &decodedJWT, nil
}
