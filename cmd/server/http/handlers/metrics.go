package handlers

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/rendering"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type MetricsHandler struct {
	*AbstractHandler
	Store  storage.Storage
	Engine rendering.Engine
}

func NewMetricsHandler(store storage.Storage, engine rendering.Engine) *MetricsHandler {
	return &MetricsHandler{
		AbstractHandler: &AbstractHandler{},
		Store:           store,
		Engine:          engine,
	}
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// URL: /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 { // 5 = "" + "update" + "metricType" + "metricName" + "metricValue"
		h.Error(w, http.StatusNotFound, "[metrics handler] bad request url, use /update/<metricType>/<metricName>/<metricValue>")
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
	default:
		possibleValues := strings.Join([]string{metrics.TypeGauge, metrics.TypeCounter}, ", ")
		h.Error(w, http.StatusNotImplemented, fmt.Sprintf("[metrics handler] bad metric type '%s', use one of: %s", metricType, possibleValues))
		return
	}

	err := h.Store.StoreMetric(metric)
	if err != nil {
		h.Error(w, http.StatusBadRequest, fmt.Sprintf("[metrics handler] error when storing metric %s: %s", metricName, err.Error()))
		return
	}

	h.Success(w, http.StatusOK, "")
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// URL: /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 { // 4 = "" + "value" + "metricType" + "metricName"
		h.Error(w, http.StatusNotFound, "[metrics handler] bad request url, use /value/<metricType>/<metricName>")
		return
	}

	metricType := parts[2]
	metricName := parts[3]

	metric, err := h.Store.GetMetric(metricType, metricName)
	if err != nil {
		var status int
		switch errors.Cause(err) {
		case storage.ErrMetricNotFound:
			status = http.StatusNotFound
		case storage.ErrMetricIncorrectType:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		h.Error(w, status, fmt.Sprintf("[metrics handler] error when getting metric from storage: %s", err.Error()))
		return
	}

	h.Success(w, http.StatusOK, metric.GetStringValue())
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	// URL: /
	mtx, err := h.Store.ListMetrics()
	if err != nil {
		h.Error(w, http.StatusInternalServerError, fmt.Sprintf("[metrics handler] error when getting metrics from storage: %s", err.Error()))
	}

	// Sort metrics by name
	sort.Slice(mtx, func(i, j int) bool {
		return strings.ToLower(mtx[i].GetName()) < strings.ToLower(mtx[j].GetName())
	})

	bytes, err := h.Engine.Render("index.html", mtx)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, fmt.Sprintf("[metrics handler] error when rendering html: %s", err.Error()))
		return
	}

	h.RenderHTML(w, http.StatusOK, bytes)
}
