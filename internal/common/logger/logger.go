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

const (
	LevelCritical = "crit"
	LevelError    = "error"
	LevelInfo     = "info"
	LevelDebug    = "debug"
)

const (
	ModeDevelopment = "dev"
	ModeProduction  = "prod"
)

type message struct {
	ctx    context.Context
	fields map[string]interface{}
}

func Init(ctx context.Context, cfg config.LoggerConfig) context.Context {
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

func New(ctx context.Context) *message {
	return &message{ctx: ctx, fields: make(map[string]interface{})}
}

func (m *message) Field(key string, value interface{}) *message {
	m.fields[key] = value
	return m
}

func (m *message) Fatalf(format string, v ...interface{}) {
	log.Ctx(m.ctx).Fatal().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Errorf(format string, v ...interface{}) {
	log.Ctx(m.ctx).Error().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Infof(format string, v ...interface{}) {
	log.Ctx(m.ctx).Info().Fields(m.fields).Msgf(format, v...)
}

func (m *message) Debugf(format string, v ...interface{}) {
	log.Ctx(m.ctx).Debug().Fields(m.fields).Msgf(format, v...)
}

func ContextFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	if requestID, ok := hlog.IDFromRequest(r); ok {
		log.Ctx(ctx).With().Bytes("request_id", requestID.Bytes())
	}
	return ctx
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
