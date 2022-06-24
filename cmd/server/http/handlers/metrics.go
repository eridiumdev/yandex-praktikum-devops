package handlers

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/routers"
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

func NewMetricsHandler(router routers.Router, store storage.Storage, engine rendering.Engine) *MetricsHandler {
	h := &MetricsHandler{
		AbstractHandler: &AbstractHandler{
			Router: router,
		},
		Store:  store,
		Engine: engine,
	}
	router.AddRoute(http.MethodGet, "/", h.List)
	router.AddRoute(http.MethodGet, "/value/{metricType}/{metricName}", h.Get)
	router.AddRoute(http.MethodPost, "/update/{metricType}/{metricName}/{metricValue}", h.Update)

	return h
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	metricType := h.Router.URLParam(r, "metricType")
	metricName := h.Router.URLParam(r, "metricName")
	metricValue := h.Router.URLParam(r, "metricValue")

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
	metricType := h.Router.URLParam(r, "metricType")
	metricName := h.Router.URLParam(r, "metricName")

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
