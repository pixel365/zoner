package metrics

import (
	"context"
	"time"
)

type Label string

const (
	ConnectionsTotal        Label = "connections_total"
	ActiveConnections       Label = "active_connections"
	FramesReadTotal         Label = "frames_read_total"
	FramesWriteTotal        Label = "frames_write_total"
	ParseErrorsTotal        Label = "parse_errors_total"
	AuthSuccessTotal        Label = "auth_success_total"
	AuthFailureTotal        Label = "auth_failure_total"
	CommandsTotal           Label = "commands_total"
	CommandsWithErrorsTotal Label = "commands_with_errors_total"
	SessionDurationMs       Label = "session_duration_ms"
)

type Collector interface {
	Inc(context.Context, Label, ...any)
	Dec(context.Context, Label, ...any)
	Duration(context.Context, Label, time.Duration, ...any)
}
