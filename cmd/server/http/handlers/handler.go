package handlers

import (
	"net/http"
)

type AbstractHandler struct {
	ActionUpdate http.HandlerFunc
}

func (h *AbstractHandler) Success(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
}

func (h *AbstractHandler) Error(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
