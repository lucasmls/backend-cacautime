package log

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

type logValues struct {
	ContextID string          `json:"contextID,omitempty"`
	Severity  infra.Severity  `json:"severity"`
	Metadata  *infra.Metadata `json:"metadata,omitempty"`
	Message   string          `json:"message"`
}

var logLevels = map[infra.Severity]int{
	infra.SeverityCritical: 1,
	infra.SeverityError:    2,
	infra.SeverityWarning:  3,
	infra.SeverityInfo:     4,
	infra.SeverityDebug:    5,
}

// ClientInput ...
type ClientInput struct {
	Level infra.Severity
	GoEnv infra.Environment
}

// Client ...
type Client struct {
	in        ClientInput
	levelCode int
}

// NewClient ...
func NewClient(in ClientInput) (*Client, *infra.Error) {
	const opName infra.OpName = "logs.NewClient"

	if in.GoEnv == "" {
		err := infra.MissingDependencyError{DependencyName: "GoEnv"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	levelCode, ok := logLevels[in.Level]
	if !ok {
		levelCode = logLevels[infra.SeverityInfo]
	}

	return &Client{
		in:        in,
		levelCode: levelCode,
	}, nil
}

func (c Client) buildLogValues(
	ctx context.Context,
	opName infra.OpName,
	level infra.Severity,
	message string,
	metadata infra.Metadata,
) string {

	values := logValues{
		Severity: level,
		Message:  fmt.Sprintf("[%s]: %s", opName, message),
	}

	if contextID := ctx.Value(infra.IDContextValueKey); contextID != nil {
		values.ContextID = fmt.Sprintf("%v", contextID)
	}

	if metadata != nil {
		values.Metadata = &metadata
	}

	logMessage, err := json.Marshal(values)

	if c.in.GoEnv == infra.EnvironmentDevelop {
		logMessage, err = json.MarshalIndent(values, "", " ")
	}

	if err != nil {
		logMessage = []byte(message)
	}

	return string(logMessage)
}

// Critical ...
func (c Client) Critical(ctx context.Context, opName infra.OpName, msg string) {
	if c.levelCode >= 1 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityCritical, msg, nil))
	}
}

// Criticalf ...
func (c Client) Criticalf(ctx context.Context, opName infra.OpName, msg string, infos ...interface{}) {
	if c.levelCode >= 1 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityCritical, fmt.Sprintf(msg, infos...), nil))
	}
}

// CriticalMetadata ...
func (c Client) CriticalMetadata(ctx context.Context, opName infra.OpName, msg string, metadata infra.Metadata) {
	if c.levelCode >= 1 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityCritical, msg, metadata))
	}
}

// Error ...
func (c Client) Error(ctx context.Context, opName infra.OpName, msg string) {
	if c.levelCode >= 2 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityError, msg, nil))
	}
}

// Errorf ...
func (c Client) Errorf(ctx context.Context, opName infra.OpName, msg string, infos ...interface{}) {
	if c.levelCode >= 2 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityError, fmt.Sprintf(msg, infos...), nil))
	}
}

// ErrorMetadata ...
func (c Client) ErrorMetadata(ctx context.Context, opName infra.OpName, msg string, metadata infra.Metadata) {
	if c.levelCode >= 2 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityError, msg, metadata))
	}
}

// Warning ...
func (c Client) Warning(ctx context.Context, opName infra.OpName, msg string) {
	if c.levelCode >= 3 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityWarning, msg, nil))
	}
}

// Warningf ...
func (c Client) Warningf(ctx context.Context, opName infra.OpName, msg string, infos ...interface{}) {
	if c.levelCode >= 3 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityWarning, fmt.Sprintf(msg, infos...), nil))
	}
}

// WarningMetadata ...
func (c Client) WarningMetadata(ctx context.Context, opName infra.OpName, msg string, metadata infra.Metadata) {
	if c.levelCode >= 3 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityWarning, msg, metadata))
	}
}

// Info ...
func (c Client) Info(ctx context.Context, opName infra.OpName, msg string) {
	if c.levelCode >= 4 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityInfo, msg, nil))
	}
}

// Infof ...
func (c Client) Infof(ctx context.Context, opName infra.OpName, msg string, infos ...interface{}) {
	if c.levelCode >= 4 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityInfo, fmt.Sprintf(msg, infos...), nil))
	}
}

// InfoMetadata ...
func (c Client) InfoMetadata(ctx context.Context, opName infra.OpName, msg string, metadata infra.Metadata) {
	if c.levelCode >= 4 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityInfo, msg, metadata))
	}
}

// Debug ...
func (c Client) Debug(ctx context.Context, opName infra.OpName, msg string) {
	if c.levelCode >= 5 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityDebug, msg, nil))
	}
}

// Debugf ...
func (c Client) Debugf(ctx context.Context, opName infra.OpName, msg string, infos ...interface{}) {
	if c.levelCode >= 5 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityDebug, fmt.Sprintf(msg, infos...), nil))
	}
}

// DebugMetadata ...
func (c Client) DebugMetadata(ctx context.Context, opName infra.OpName, msg string, metadata infra.Metadata) {
	if c.levelCode >= 5 {
		fmt.Println(c.buildLogValues(ctx, opName, infra.SeverityDebug, msg, metadata))
	}
}
