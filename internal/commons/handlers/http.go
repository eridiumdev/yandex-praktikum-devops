package handlers

import (
	"encoding/json"
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/routing"
)

type HTTPHandler struct {
	Router routing.Router
}

func (h *HTTPHandler) PlainText(w http.ResponseWriter, status int, body string) {
	h.write(w, status, []byte(body), "text/plain; charset=utf-8")
}

func (h *HTTPHandler) HTML(w http.ResponseWriter, body []byte) {
	h.write(w, http.StatusOK, body, "text/html; charset=utf-8")
}

func (h *HTTPHandler) JSON(w http.ResponseWriter, status int, data struct{}) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("error when marshaling data %v, responding with an empty json struct", data)
		body = []byte(`{}`)
		status = http.StatusInternalServerError
	}
	h.write(w, status, body, "application/json; charset=utf-8")
}

func (h *HTTPHandler) write(w http.ResponseWriter, status int, body []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if body != nil {
		_, err := w.Write(body)
		if err != nil {
			logger.Errorf("could not write body to writer")
		}
	}
}
