package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/handlers"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

const (
	ErrStringInvalidJSON       = "invalid JSON"
	ErrStringInvalidMetricType = "invalid metric type"
	ErrStringInvalidHash       = "invalid hash"
	ErrStringMetricNotFound    = "metric not found"
	ErrStringRenderingError    = "rendering error"
)

type MetricsHandler struct {
	*handlers.HTTPHandler
	service  MetricsService
	renderer MetricsRenderer
	factory  MetricsRequestResponseFactory
	hasher   MetricsHasher
}

func NewMetricsHandler(
	service MetricsService,
	renderer MetricsRenderer,
	factory MetricsRequestResponseFactory,
	hasher MetricsHasher,
) *MetricsHandler {
	return &MetricsHandler{
		HTTPHandler: &handlers.HTTPHandler{},
		service:     service,
		renderer:    renderer,
		factory:     factory,
		hasher:      hasher,
	}
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)
	var req domain.UpdateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.New(ctx).Errorf("[metrics handler] received invalid JSON: %s", err.Error())
		h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidJSON)
		return
	}

	if !domain.IsValidMetricType(req.MType) {
		logger.New(ctx).Errorf("[metrics handler] received invalid req type '%s'", req.MType)
		h.PlainText(ctx, w, http.StatusNotImplemented, ErrStringInvalidMetricType)
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
	// Validate hash
	if req.Hash != "" && !h.hasher.Check(ctx, metric, req.Hash) {
		logger.New(ctx).Errorf("[metrics handler] provided hash is invalid")
		h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidHash)
		return
	}
	updatedMetric, _ := h.service.Update(metric)

	h.JSON(ctx, w, http.StatusOK, h.factory.BuildUpdateMetricResponse(ctx, updatedMetric))
}

func (h *MetricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)
	var req domain.GetMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.New(ctx).Errorf("[metrics handler] received invalid JSON: %s", err.Error())
		h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidJSON)
		return
	}

	if !domain.IsValidMetricType(req.MType) {
		logger.New(ctx).Errorf("[metrics handler] received invalid metric type '%s'", req.MType)
		h.PlainText(ctx, w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}

	metric, found := h.service.Get(req.ID)
	if !found || metric.Type != req.MType {
		logger.New(ctx).Errorf("[metrics handler] metric '%s/%s' not found", req.MType, req.ID)
		h.PlainText(ctx, w, http.StatusNotFound, ErrStringMetricNotFound)
		return
	}

	h.JSON(ctx, w, http.StatusOK, h.factory.BuildGetMetricResponse(ctx, metric))
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)
	list := h.service.List()

	// Sort metrics by name
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	html, err := h.renderer.RenderList(list)
	if err != nil {
		logger.New(ctx).Errorf(fmt.Sprintf("[metrics handler] error when rendering html: %s", err.Error()))
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringRenderingError)
		return
	}

	h.HTML(ctx, w, html)
}
