package infra

import (
	"context"
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

// ErrorKind ...
type ErrorKind int

const (
	// KindBadRequest ...
	KindBadRequest ErrorKind = http.StatusBadRequest
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
