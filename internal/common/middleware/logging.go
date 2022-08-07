package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

// BaseLoggingMiddleware is required to be included in order for other logging middlewares to work
func BaseLoggingMiddleware(next http.Handler) http.Handler {
	return hlog.NewHandler(log.Logger)(next)
}

func AddRequestID(next http.Handler) http.Handler {
	return hlog.RequestIDHandler("request_id", "X-Request-Id")(next)
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.New(r.Context()).
			Field("host", r.Host).
			Field("method", r.Method).
			Field("url", r.URL.String()).
			Field("content_type", r.Header.Get("Content-Type")).
			Field("content_encoding", r.Header.Get("Content-Encoding")).
			Infof("--> HTTP request")
		next.ServeHTTP(w, r)
	})
}

func LogResponses(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		defer func() {
			logger.New(r.Context()).
				Field("status", ww.Status()).
				Field("content_type", ww.Header().Get("Content-Type")).
				Field("duration_ms", time.Since(start).Round(time.Microsecond)).
				Infof("<-- HTTP response")
		}()
		next.ServeHTTP(ww, r)
	})
}

func LogRequestsWithBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var body []byte
		if r.Body != nil {
			tee := io.TeeReader(r.Body, &buf)
			body, _ = ioutil.ReadAll(tee)
			r.Body = ioutil.NopCloser(&buf)
		}
		logger.New(r.Context()).
			Field("host", r.Host).
			Field("method", r.Method).
			Field("url", r.URL.String()).
			Field("content_type", r.Header.Get("Content-Type")).
			Field("content_encoding", r.Header.Get("Content-Encoding")).
			Field("body", string(body)).
			Infof("--> HTTP request")
		next.ServeHTTP(w, r)
	})
}

func LogResponsesWithBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ww.Tee(&buf)
		start := time.Now()
		defer func() {
			body, _ := ioutil.ReadAll(&buf)
			logger.New(r.Context()).
				Field("status", ww.Status()).
				Field("content_type", ww.Header().Get("Content-Type")).
				Field("body", string(body)).
				Field("duration_ms", time.Since(start).Round(time.Microsecond)).
				Infof("<-- HTTP response")
		}()

		next.ServeHTTP(ww, r)
	})
}
