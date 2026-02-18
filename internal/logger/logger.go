package logger

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, error)
	Logf(string, ...any)
	Component(string) Logger
	Func(string) Logger
	ClientId(string) Logger
	SessionId(string) Logger
}
