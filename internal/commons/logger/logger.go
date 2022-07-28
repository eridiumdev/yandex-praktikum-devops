package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

const (
	LevelCritical = iota
	LevelError
	LevelInfo
	LevelDebug
)

const (
	ModeDevelopment = iota
	ModeProduction
)

func Init(level uint8, mode uint8) {
	zerolog.SetGlobalLevel(convertToZerologLevel(level))

	switch mode {
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
}

func Middleware(next http.Handler) http.Handler {
	middleware1 := hlog.NewHandler(log.Logger)
	middleware2 := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("http request")
	})
	return middleware1(middleware2(next))
}

func Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func convertToZerologLevel(level uint8) zerolog.Level {
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
