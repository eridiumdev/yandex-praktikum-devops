package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type HTTPHandler struct{}

func (h *HTTPHandler) PlainText(ctx context.Context, w http.ResponseWriter, status int, body string) {
	h.write(ctx, w, status, []byte(body), "text/plain; charset=utf-8")
}

func (h *HTTPHandler) HTML(ctx context.Context, w http.ResponseWriter, body []byte) {
	h.write(ctx, w, http.StatusOK, body, "text/html; charset=utf-8")
}

func (h *HTTPHandler) JSON(ctx context.Context, w http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.New(ctx).Errorf("error when marshaling data %v, responding with an empty json struct", data)
		body = []byte(`{}`)
		status = http.StatusInternalServerError
	}
	h.write(ctx, w, status, body, "application/json; charset=utf-8")
}

func (h *HTTPHandler) write(ctx context.Context, w http.ResponseWriter, status int, body []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if body != nil {
		_, err := w.Write(body)
		if err != nil {
			logger.New(ctx).Errorf("could not write body to writer")
		}
	}
}
