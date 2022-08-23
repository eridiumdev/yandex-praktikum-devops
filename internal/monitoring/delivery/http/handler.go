package http

import (
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/handlers"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type MonitoringHandler struct {
	*handlers.HTTPHandler
	components []Pingable
}

func NewMonitoringHandler(components ...Pingable) *MonitoringHandler {
	return &MonitoringHandler{
		HTTPHandler: &handlers.HTTPHandler{},
		components:  components,
	}
}

func (h *MonitoringHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)

	statusOk := true

	for _, component := range h.components {
		statusOk = statusOk && component.Ping(ctx)
	}

	if !statusOk {
		h.PlainText(ctx, w, http.StatusInternalServerError, "")
		return
	}

	h.PlainText(ctx, w, http.StatusOK, "")
}
