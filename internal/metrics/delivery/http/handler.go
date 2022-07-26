package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/handlers"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

const (
	ErrStringInvalidJSON       = "invalid JSON"
	ErrStringInvalidMetricType = "invalid metric type"
	ErrStringMetricNotFound    = "metric not found"
	ErrStringRenderingError    = "rendering error"
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
	router.AddRoute(http.MethodPost, "/value", h.Get)
	router.AddRoute(http.MethodPost, "/update", h.Update)

	return h
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("[metrics handler] received invalid JSON: %s", err.Error())
		h.PlainText(w, http.StatusBadRequest, ErrStringInvalidJSON)
		return
	}

	if !domain.IsValidMetricType(req.MType) {
		logger.Errorf("[metrics handler] received invalid req type '%s'", req.MType)
		h.PlainText(w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	metric := domain.Metric{
		Name: req.ID,
		Type: req.MType,
	}
	if req.Delta != nil {
		metric.Counter = domain.Counter(*req.Delta)
	}
	if req.Value != nil {
		metric.Gauge = domain.Gauge(*req.Value)
	}
	h.service.Update(metric)

	h.JSON(w, http.StatusOK, domain.PrepareUpdateMetricResponse(metric))
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	var req domain.GetMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("[metrics handler] received invalid JSON: %s", err.Error())
		h.PlainText(w, http.StatusBadRequest, ErrStringInvalidJSON)
		return
	}

	if !domain.IsValidMetricType(req.MType) {
		logger.Errorf("[metrics handler] received invalid metric type '%s'", req.MType)
		h.PlainText(w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	metric, found := h.service.Get(req.ID)
	if !found || metric.Type != req.MType {
		logger.Errorf("[metrics handler] metric '%s/%s' not found", req.MType, req.ID)
		h.PlainText(w, http.StatusNotFound, ErrStringMetricNotFound)
		return
	}

	h.JSON(w, http.StatusOK, domain.PrepareGetMetricResponse(metric))
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	list := h.service.List()

	// Sort metrics by name
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	html, err := h.renderer.RenderList(list)
	if err != nil {
		logger.Errorf(fmt.Sprintf("[metrics handler] error when rendering html: %s", err.Error()))
		h.PlainText(w, http.StatusInternalServerError, ErrStringRenderingError)
		return
	}

	h.HTML(w, html)
}
