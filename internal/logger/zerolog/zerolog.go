package zerolog

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/pixel365/zoner/internal/logger"
)

type Log struct {
	logger zerolog.Logger
}

func (o *Log) Debug(msg string, params ...any) {
	if len(params) > 0 {
		o.logger.Debug().Msg(fmt.Sprintf(msg, params...))
		return
	}
	o.logger.Debug().Msg(msg)
}

func (o *Log) Info(msg string, params ...any) {
	if len(params) > 0 {
		o.logger.Info().Msg(fmt.Sprintf(msg, params...))
		return
	}
	o.logger.Info().Msg(msg)
}

func (o *Log) Warn(msg string, params ...any) {
	if len(params) > 0 {
		o.logger.Warn().Msg(fmt.Sprintf(msg, params...))
		return
	}
	o.logger.Warn().Msg(msg)
}

func (o *Log) Error(msg string, err error) {
	if err == nil {
		o.logger.Error().Msg(msg)
		return
	}
	o.logger.Error().
		Str("error.message", err.Error()).
		Str("error.type", fmt.Sprintf("%T", err)).
		Msg(msg)
}

func (o *Log) Logf(format string, args ...any) {
	o.logger.Info().Msgf(format, args...)
}

func (o *Log) Component(name string) logger.Logger {
	if name == "" {
		return o
	}

	l := o.logger.With().
		Str("log.logger", name).
		Logger()

	return &Log{logger: l}
}

func (o *Log) Func(name string) logger.Logger {
	if name == "" {
		return o
	}

	l := o.logger.With().
		Str("log.origin.function", name).
		Logger()

	return &Log{logger: l}
}

func (o *Log) WithUserId(id string) logger.Logger {
	if id == "" {
		return o
	}

	l := o.logger.With().
		Str("user.id", id).
		Logger()

	return &Log{logger: l}
}

func (o *Log) WithSessionId(id string) logger.Logger {
	if id == "" {
		return o
	}

	l := o.logger.With().
		Str("session.id", id).
		Logger()

	return &Log{logger: l}
}

func (o *Log) WithAddress(addr string) logger.Logger {
	if addr == "" {
		return o
	}

	l := o.logger.With().
		Str("client.address", addr).
		Logger()

	return &Log{logger: l}
}

func (o *Log) WithEventDuration(duration time.Duration) logger.Logger {
	if duration <= 0 {
		return o
	}

	l := o.logger.With().
		Int64("event.duration_ms", duration.Milliseconds()).
		Logger()

	return &Log{logger: l}
}

func NewLogger(cfg *logger.Config, writers ...io.Writer) *Log {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "log.level"

	var outs []io.Writer
	if len(writers) == 0 {
		outs = []io.Writer{os.Stdout}
	} else {
		outs = writers
	}

	multi := zerolog.MultiLevelWriter(outs...)
	base := zerolog.New(multi).
		Level(zerolog.Level(cfg.LogLevel)).
		With().
		Timestamp().
		Str("ecs.version", cfg.ECSVersion).
		Str("service.name", cfg.Service.Name).
		Str("service.version", cfg.Service.Version).
		Str("service.environment", cfg.Service.Environment).
		Str("service.instance.id", cfg.Service.InstanceID).
		Str("event.dataset", cfg.EventDataset).
		Int("process.pid", os.Getpid()).
		Str("host.name", hostname()).
		Logger()

	if cfg.LoggerName != "" {
		base = base.With().Str("log.logger", cfg.LoggerName).Logger()
	}

	return &Log{logger: base}
}
