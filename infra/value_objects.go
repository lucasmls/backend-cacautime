package infra

import (
	"context"
	"fmt"
	"net/http"
)

// OpName ...
type OpName string

// Severity ...
type Severity string

const (
	// SeverityCritical ...
	SeverityCritical Severity = "Critical"
	// SeverityError ...
	SeverityError Severity = "Error"
	// SeverityWarning ...
	SeverityWarning Severity = "Warning"
	// SeverityInfo ...
	SeverityInfo Severity = "Info"
	// SeverityDebug ...
	SeverityDebug Severity = "Debug"
)

// Environment ...
type Environment string

const (
	// EnvironmentDevelop ...
	EnvironmentDevelop Environment = "development"
	// EnvironmentProduction ...
	EnvironmentProduction Environment = "production"
	// EnvironmentStaging ...
	EnvironmentStaging Environment = "staging"
)

const (
	// IDContextValueKey ...
	IDContextValueKey string = "contextID"
)

// ErrorKind ...
type ErrorKind int

const (
	// KindBadRequest ...
	KindBadRequest ErrorKind = http.StatusBadRequest
	// KindUnauthorized ...
	KindUnauthorized ErrorKind = http.StatusUnauthorized
	// KindNotFound ...
	KindNotFound ErrorKind = http.StatusNotFound
	// KindUnexpected ...
	KindUnexpected ErrorKind = http.StatusInternalServerError
	// KindExpected ...
	KindExpected ErrorKind = http.StatusOK
)

// Metadata ...
type Metadata map[string]interface{}

// Merge ...
func (md Metadata) Merge(newMd *Metadata) Metadata {
	if newMd == nil {
		return md
	}

	for key, value := range *newMd {
		md[key] = value
	}

	return md
}

// Error ...
type Error struct {
	Ctx      context.Context `json:"-"`
	Err      error           `json:"-"`
	Severity Severity        `json:"severity"`
	OpName   OpName          `json:"opName"`
	Kind     ErrorKind       `json:"kind"`
	Metadata Metadata        `json:"metadata"`
}

// Error ...
func (e Error) Error() string {
	return e.Err.Error()
}

// MissingEnvironmentVariableError ...
type MissingEnvironmentVariableError struct {
	EnvVarName string
}

func (e MissingEnvironmentVariableError) Error() string {
	return fmt.Sprintf("missing required environment variable: %s", e.EnvVarName)
}

// MissingDependencyError ...
type MissingDependencyError struct {
	DependencyName string
}

func (e MissingDependencyError) Error() string {
	return fmt.Sprintf("missing required dependency: %s", e.DependencyName)
}

// MinimumValueError ...
type MinimumValueError struct {
	EnvVarName      string
	MinimumRequired int
}

func (e MinimumValueError) Error() string {
	return fmt.Sprintf("missing value: %s - minimum required: %d", e.EnvVarName, e.MinimumRequired)
}

// DecodedJWT ...
type DecodedJWT struct {
	UserID string `json:"userId"`
	Exp    int64  `json:"exp"`
}
