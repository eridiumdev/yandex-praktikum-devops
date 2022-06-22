package handlers

import (
	"net/http"
)

type AbstractHandler struct{}

func (h *AbstractHandler) Success(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	if body != "" {
		w.Write([]byte(body))
	}
}

func (h *AbstractHandler) RenderHTML(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		w.Write(body)
	}
}

func (h *AbstractHandler) Error(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
