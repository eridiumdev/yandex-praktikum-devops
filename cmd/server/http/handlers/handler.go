package handlers

import (
	"net/http"
)

type AbstractHandler struct {
	ActionUpdate http.HandlerFunc
}

func (h *AbstractHandler) Success(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
}

func (h *AbstractHandler) Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, msg, status)
}
