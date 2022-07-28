package http

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/handlers"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

const (
	ErrStringInvalidMetricType  = "invalid metric type"
	ErrStringInvalidMetricValue = "invalid metric value"
	ErrStringMetricNotFound     = "metric not found"
	ErrStringUpdateError        = "update error"
	ErrStringRetrieveError      = "retrieve error"
	ErrStringRenderingError     = "rendering error"
)

type MetricsHandler struct {
	*handlers.HTTPHandler
	service  MetricsService
	renderer MetricsRenderer
}

func NewMetricsHandler(router routing.Router, service MetricsService, renderer MetricsRenderer) *MetricsHandler {
	h := &MetricsHandler{
		HTTPHandler: &handlers.HTTPHandler{
			Router: router,
		},
		service:  service,
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
	metricValueRaw := h.Router.URLParam(r, "metricValue")

	var updateErr error

	switch metricType {
	case domain.TypeCounter:
		val, err := strconv.ParseInt(metricValueRaw, 10, 64)
		if err != nil {
			logger.Errorf("[metrics handler] received invalid metric value '%s': %s", metricValueRaw, err.Error())
			h.PlainText(w, http.StatusBadRequest, ErrStringInvalidMetricValue)
			return
		}
		_, updateErr = h.service.UpdateCounter(metricName, domain.Counter(val))
	case domain.TypeGauge:
		val, err := strconv.ParseFloat(metricValueRaw, 64)
		if err != nil {
			logger.Errorf("[metrics handler] received invalid metric value '%s': %s", metricValueRaw, err.Error())
			h.PlainText(w, http.StatusBadRequest, ErrStringInvalidMetricValue)
			return
		}
		_, updateErr = h.service.UpdateGauge(metricName, domain.Gauge(val))
	default:
		logger.Errorf("[metrics handler] received invalid metric type '%s'", metricType)
		h.PlainText(w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	if updateErr != nil {
		logger.Errorf("[metrics handler] error when updating metric %s: %s", metricName, updateErr.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringUpdateError)
		return
	}

	h.PlainText(w, http.StatusOK, "")
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	metricType := h.Router.URLParam(r, "metricType")
	metricName := h.Router.URLParam(r, "metricName")

	if !domain.IsValidMetricType(metricType) {
		logger.Errorf("[metrics handler] received invalid metric type '%s'", metricType)
		h.PlainText(w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	metric, err := h.service.Get(metricName)
	if err != nil {
		logger.Errorf("[metrics handler] error when retrieving metric: %s", err.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringRetrieveError)
		return
	}
	if metric == nil || metric.Type() != metricType {
		logger.Errorf("[metrics handler] metric '%s/%s' not found", metricType, metricName)
		h.PlainText(w, http.StatusNotFound, ErrStringMetricNotFound)
		return
	}

	h.PlainText(w, http.StatusOK, metric.StringValue())
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List()
	if err != nil {
		logger.Errorf("[metrics handler] error when retrieving metric list: %s", err.Error())
		h.PlainText(w, http.StatusInternalServerError, ErrStringRetrieveError)
		return
	}

	// Sort metrics by name
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name()) < strings.ToLower(list[j].Name())
	})

	html, err := h.renderer.RenderList(list)
	if err != nil {
		logger.Errorf(fmt.Sprintf("[metrics handler] error when rendering html: %s", err.Error()))
		h.PlainText(w, http.StatusInternalServerError, ErrStringRenderingError)
		return
	}

	h.HTML(w, html)
}
