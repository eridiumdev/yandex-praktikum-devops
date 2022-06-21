package handlers

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MetricsHandler struct {
	*AbstractHandler
	Store storage.Storage
}

func NewMetricsHandler(store storage.Storage) *MetricsHandler {
	return &MetricsHandler{
		AbstractHandler: &AbstractHandler{},
		Store: store,
	}
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// URL: /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 { // 5 = "" + "update" + "metricType" + "metricName" + "metricValue"
		h.Error(w, http.StatusBadRequest, "[metrics handler] bad request url, use /update/<metricType>/<metricName>/<metricValue>")
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	var metric metrics.Metric

	switch metricType {
	case metrics.TypeCounter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			h.Error(w, http.StatusBadRequest, fmt.Sprintf("[metrics handler] bad metric value '%s': %s", metricValue, err.Error()))
			return
		}
		metric = metrics.CounterMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metricName,
			},
			Value: metrics.Counter(val),
		}
	case metrics.TypeGauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			h.Error(w, http.StatusBadRequest, fmt.Sprintf("[metrics handler] bad metric value '%s': %s", metricValue, err.Error()))
			return
		}
		metric = metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metricName,
			},
			Value: metrics.Gauge(val),
		}
	}

	err := h.Store.StoreMetric(metric)
	if err != nil {
		h.Error(w, http.StatusBadRequest, fmt.Sprintf("[metrics handler] error when storing metric %s: %s", metricName, err.Error()))
		return
	}

	h.Success(w, http.StatusOK)
}
