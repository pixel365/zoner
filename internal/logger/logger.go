package logger

import "time"

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, error)
	Logf(string, ...any)
	Component(string) Logger
	Func(string) Logger
	WithUsername(string) Logger
	WithUserId(int64) Logger
	WithSessionId(string) Logger
	WithAddress(string) Logger
	WithEventDuration(time.Duration) Logger
}
