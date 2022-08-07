package middleware

import (
	"compress/gzip"
	"net/http"

	chi "github.com/go-chi/chi/middleware"
)

func DecompressRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}
		next.ServeHTTP(w, r)
	})
}

func CompressResponses(next http.Handler) http.Handler {
	return chi.Compress(5)(next)
}
