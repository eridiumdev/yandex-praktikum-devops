package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	_http "eridiumdev/yandex-praktikum-go-devops/cmd/server/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/server"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/server/rendering"
)

const (
	ErrStringInvalidMetricType  = "invalid metric type"
	ErrStringInvalidMetricValue = "invalid metric value"
	ErrStringMetricNotFound     = "metric not found"
	ErrStringDatabaseError      = "database error"
	ErrStringRenderingError     = "rendering error"
)

type MetricsHandler struct {
	*AbstractHandler
	repo     server.MetricsRepository
	renderer server.MetricsRenderer
}

func NewMetricsHandler(router _http.Router, repo server.MetricsRepository, renderer server.MetricsRenderer) *MetricsHandler {
	h := &MetricsHandler{
		AbstractHandler: &AbstractHandler{
			Router: router,
		},
		repo:     repo,
		renderer: renderer,
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

	var metric domain.Metric

	switch metricType {
	case domain.TypeCounter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Errorf("[metrics handler] received invalid metric value '%s': %s", metricValue, err.Error())
			h.PlainText(w, http.StatusBadRequest, ErrStringInvalidMetricValue)
			return
		}
		metric = domain.NewCounter(metricName, domain.Counter(val))
	case domain.TypeGauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Errorf("[metrics handler] received invalid metric value '%s': %s", metricValue, err.Error())
			h.PlainText(w, http.StatusBadRequest, ErrStringInvalidMetricValue)
			return
		}
		metric = domain.NewGauge(metricName, domain.Gauge(val))
	default:
		logger.Errorf("[metrics handler] received invalid metric type '%s'", metricType)
		h.PlainText(w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	err := h.repo.Store(metric)
	if err != nil {
		logger.Errorf("[metrics handler] error when storing metric %s: %s", metricName, err.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}

	h.PlainText(w, http.StatusOK, "")
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	metricType := h.Router.URLParam(r, "metricType")
	metricName := h.Router.URLParam(r, "metricName")

	metric, err := h.repo.Get(metricName)
	if err != nil {
		logger.Errorf("[metrics handler] error when getting metric from repository: %s", err.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}
	if metric == nil {
		logger.Errorf("[metrics handler] metric '%s' not found", metricName)
		h.PlainText(w, http.StatusNotFound, ErrStringMetricNotFound)
		return
	}
	if metric.Type() != metricType {
		logger.Errorf("[metrics handler] invalid metric type requested '%s', '%s' wanted", metricType, metric.Type())
		h.PlainText(w, http.StatusBadRequest, ErrStringInvalidMetricType)
		return
	}

	h.PlainText(w, http.StatusOK, metric.StringValue())
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	mtx, err := h.repo.List()
	if err != nil {
		logger.Errorf("[metrics handler] error when getting metrics from repository: %s", err.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringDatabaseError)
	}

	// Sort metrics by name
	sort.Slice(mtx, func(i, j int) bool {
		return strings.ToLower(mtx[i].Name()) < strings.ToLower(mtx[j].Name())
	})

	html, err := h.renderer.RenderList(rendering.MetricsListTemplate, mtx)
	if err != nil {
		logger.Errorf(fmt.Sprintf("[metrics handler] error when rendering html: %s", err.Error()))
		h.PlainText(w, http.StatusInternalServerError, ErrStringRenderingError)
		return
	}

	h.HTML(w, html)
}
