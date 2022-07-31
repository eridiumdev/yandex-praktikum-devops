package logger

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
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

type message struct {
	ctx    context.Context
	fields map[string]interface{}
}

func Init(ctx context.Context, level uint8, mode uint8) context.Context {
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
	return log.Logger.WithContext(ctx)
}

func Middleware(next http.Handler) http.Handler {
	middlewareChain := make([]func(http.Handler) http.Handler, 0)

	middlewareChain = append(middlewareChain, LogResponses)
	middlewareChain = append(middlewareChain, LogRequests)
	middlewareChain = append(middlewareChain, hlog.RequestIDHandler("request_id", "X-Request-Id"))
	middlewareChain = append(middlewareChain, hlog.NewHandler(log.Logger))

	for _, mdlw := range middlewareChain {
		next = mdlw(next)
	}
	return next
}

func ContextFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	if requestID, ok := hlog.IDFromRequest(r); ok {
		log.Ctx(ctx).With().Bytes("request_id", requestID.Bytes())
	}
	return ctx
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var body []byte
		if r.Body != nil {
			tee := io.TeeReader(r.Body, &buf)
			body, _ = ioutil.ReadAll(tee)
			r.Body = ioutil.NopCloser(&buf)
		}
		New(r.Context()).
			Field("host", r.Host).
			Field("method", r.Method).
			Field("url", r.URL.String()).
			Field("content_type", r.Header.Get("Content-Type")).
			Field("body", string(body)).
			Infof("--> HTTP request")
		next.ServeHTTP(w, r)
	})
}

func LogResponses(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ww.Tee(&buf)
		start := time.Now()
		defer func() {
			body, _ := ioutil.ReadAll(&buf)
			New(r.Context()).
				Field("status", ww.Status()).
				Field("content_type", ww.Header().Get("Content-Type")).
				Field("body", string(body)).
				Field("duration_ms", time.Since(start).Round(time.Microsecond)).
				Infof("<-- HTTP response")
		}()

		next.ServeHTTP(ww, r)
	})
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
