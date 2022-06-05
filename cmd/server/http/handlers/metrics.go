package handlers

import (
	"net/http"
	"strings"
)

type MetricsHandler struct {
	*AbstractHandler
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		AbstractHandler: &AbstractHandler{},
	}
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// URL: /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 { // 5 = "" + "update" + "metricType" + "metricName" + "metricValue"
		h.Error(w, http.StatusBadRequest, "bad request url, use /update/<metricType>/<metricName>/<metricValue>")
		return
	}

	h.Success(w, http.StatusOK)
}
