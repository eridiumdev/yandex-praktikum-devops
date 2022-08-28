package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"eridiumdev/yandex-praktikum-go-devops/config"
)

func InitZerolog(ctx context.Context, cfg config.LoggerConfig) context.Context {
	zerolog.SetGlobalLevel(convertToZerologLevel(cfg.Level))

	switch cfg.Mode {
	case ModeProduction:
		log.Logger = log.Output(os.Stdout)
	case ModeDevelopment:
		fallthrough
	default:
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
	return log.Logger.WithContext(ctx)
}

func ContextFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	if requestID, ok := hlog.IDFromRequest(r); ok {
		log.Ctx(ctx).With().Bytes("request_id", requestID.Bytes())
	}
	return ctx
}

func (m *message) Fatalf(format string, v ...interface{}) {
	m.getZerologLogger(m.ctx).Fatal().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Errorf(format string, v ...interface{}) {
	m.getZerologLogger(m.ctx).Error().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Warnf(format string, v ...interface{}) {
	m.getZerologLogger(m.ctx).Warn().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Infof(format string, v ...interface{}) {
	m.getZerologLogger(m.ctx).Info().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Debugf(format string, v ...interface{}) {
	m.getZerologLogger(m.ctx).Debug().Fields(m.fields).Msgf(format, v...)
}

func (m *message) getZerologLogger(ctx context.Context) *zerolog.Logger {
	logger := log.Ctx(ctx)
	if logger.GetLevel() == zerolog.Disabled {
		return &log.Logger
	}
	return logger
}

func convertToZerologLevel(level string) zerolog.Level {
	switch level {
	case LevelCritical:
		return zerolog.FatalLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		fallthrough
	default:
		return zerolog.InfoLevel
	}
}
