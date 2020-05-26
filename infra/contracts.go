package infra

import (
	"context"
	"database/sql/driver"
)

// ObjectID represents a document identifier
type ObjectID int

// LogProvider ...
type LogProvider interface {
	Critical(context.Context, OpName, string)
	Criticalf(context.Context, OpName, string, ...interface{})
	CriticalMetadata(context.Context, OpName, string, Metadata)
	Error(context.Context, OpName, string)
	Errorf(context.Context, OpName, string, ...interface{})
	ErrorMetadata(context.Context, OpName, string, Metadata)
	Warning(context.Context, OpName, string)
	Warningf(context.Context, OpName, string, ...interface{})
	WarningMetadata(context.Context, OpName, string, Metadata)
	Info(context.Context, OpName, string)
	Infof(context.Context, OpName, string, ...interface{})
	InfoMetadata(context.Context, OpName, string, Metadata)
	Debug(context.Context, OpName, string)
	Debugf(context.Context, OpName, string, ...interface{})
	DebugMetadata(context.Context, OpName, string, Metadata)
}

// RelationalDatabaseProvider ...
type RelationalDatabaseProvider interface {
	Query(context.Context, string, ...interface{}) Decoder
	QueryAll(context.Context, string, ...interface{}) (Cursor, *Error)
	Execute(context.Context, string, ...interface{}) (driver.Result, *Error)
}

// Entity represents an abstraction of an entity in database
type Entity interface{}

// Cursor ...
type Cursor interface {
	Next(context.Context) bool
	Decode(context.Context, Entity) *Error
	Close(context.Context) *Error
}

// Decoder ...
type Decoder interface {
	Decode(context.Context, Entity) *Error
}
