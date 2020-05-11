package errors

import (
	"context"
	"errors"

	"github.com/lucasmls/backend-cacautime/infra"
)

// New ...
func New(args ...interface{}) *infra.Error {
	err := infra.Error{
		Metadata: infra.Metadata{},
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case context.Context:
			err.Ctx = arg
		case error:
			err.Err = arg
		case infra.OpName:
			err.OpName = arg
		case infra.ErrorKind:
			err.Kind = arg
		case infra.Severity:
			err.Severity = arg
		case infra.Metadata:
			err.Metadata = err.Metadata.Merge(&arg)
		case string:
			err.Err = errors.New(arg)
		}
	}

	return &err
}
