package metrics

import (
	"context"
	"time"
)

type Label string

func (l Label) Description() string {
	switch l {
	case ConnectionsTotal:
		return "Total number of accepted client connections."
	case ParseErrorsTotal:
		return "Total number of EPP command parse errors."
	case AuthSuccessTotal:
		return "Total number of successful login attempts."
	case AuthFailureTotal:
		return "Total number of failed login attempts."
	case CommandsTotal:
		return "Total number of parsed EPP commands."
	case RequestsTotal:
		return "Total number of inbound network requests accepted by listener."
	case CommandsWithErrorsTotal:
		return "Total number of command handling errors."
	case ActiveConnections:
		return "Current number of active client connections."
	case FramesReadTotal:
		return "Distribution of inbound EPP frame payload sizes in bytes."
	case FramesWriteTotal:
		return "Distribution of outbound EPP frame payload sizes in bytes."
	case SessionDurationMs:
		return "Distribution of client session durations in milliseconds."
	default:
		return "Unknown label: " + string(l)
	}
}

const (
	// Int64Counter
	ConnectionsTotal        Label = "connections_total"
	ParseErrorsTotal        Label = "parse_errors_total"
	AuthSuccessTotal        Label = "auth_success_total"
	AuthFailureTotal        Label = "auth_failure_total"
	CommandsTotal           Label = "commands_total"
	RequestsTotal           Label = "requests_total"
	CommandsWithErrorsTotal Label = "commands_with_errors_total"

	// Int64UpDownCounter
	ActiveConnections Label = "active_connections"

	// Int64Histogram
	FramesReadTotal  Label = "frames_read_total"
	FramesWriteTotal Label = "frames_write_total"

	// Float64Histogram
	SessionDurationMs Label = "session_duration_ms"
)

type IncBytesFunc func(context.Context, Label, int64, ...any)

type Collector interface {
	Inc(context.Context, Label, ...any)
	Dec(context.Context, Label, ...any)
	Duration(context.Context, Label, time.Duration, ...any)
	IncBytes(context.Context, Label, int64, ...any)
	Shutdown(context.Context) error
}
