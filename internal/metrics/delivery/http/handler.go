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
	ErrStringDatabaseError     = "database error"
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
		logger.New(ctx).Errorf("[metrics handler] received invalid metric type '%s'", req.MType)
		h.PlainText(ctx, w, http.StatusNotImplemented, ErrStringInvalidMetricType)
		return
	}
	metric := req.TranslateToMetric()

	// Validate hash
	if req.Hash != "" && !h.hasher.Check(ctx, metric, req.Hash) {
		logger.New(ctx).Errorf("[metrics handler] provided hash is invalid")
		h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidHash)
		return
	}

	updatedMetric, err := h.service.Update(ctx, metric)
	if err != nil {
		logger.New(ctx).Errorf("[metrics handler] error when updating metric: %s", err.Error())
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}

	h.JSON(ctx, w, http.StatusOK, h.factory.BuildUpdateMetricResponse(ctx, updatedMetric))
}

func (h *MetricsHandler) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)
	var req []domain.UpdateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.New(ctx).Errorf("[metrics handler] received invalid JSON: %s", err.Error())
		h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidJSON)
		return
	}

	metrics := make([]domain.Metric, 0)
	for _, reqMetric := range req {
		if !domain.IsValidMetricType(reqMetric.MType) {
			logger.New(ctx).Errorf("[metrics handler] received invalid metric type '%s'", reqMetric.MType)
			h.PlainText(ctx, w, http.StatusNotImplemented, ErrStringInvalidMetricType)
			return
		}
		metric := reqMetric.TranslateToMetric()

		// Validate hash
		if reqMetric.Hash != "" && !h.hasher.Check(ctx, metric, reqMetric.Hash) {
			logger.New(ctx).Errorf("[metrics handler] provided hash for metric %s is invalid", reqMetric.ID)
			h.PlainText(ctx, w, http.StatusBadRequest, ErrStringInvalidHash)
			return
		}
		metrics = append(metrics, metric)
	}

	updatedMetrics, err := h.service.UpdateMany(ctx, metrics)
	if err != nil {
		logger.New(ctx).Errorf("[metrics handler] error when updating metric: %s", err.Error())
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}

	body, err := json.Marshal(h.factory.BuildUpdateBatchMetricResponse(ctx, updatedMetrics))
	if err != nil {
		logger.New(ctx).Errorf("[metrics handler] error when marshaling updated metrics: %s", err.Error())
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringRenderingError)
		return
	}
	// Plaintext instead of JSON to hack Yandex-practicum tests (does not work with JSON-array)
	h.PlainText(ctx, w, http.StatusOK, string(body))
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

	metric, found, err := h.service.Get(ctx, req.ID)
	if err != nil {
		logger.New(ctx).Errorf("[metrics handler] error when getting metric: %s", err.Error())
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}
	if !found || metric.Type != req.MType {
		logger.New(ctx).Errorf("[metrics handler] metric '%s/%s' not found", req.MType, req.ID)
		h.PlainText(ctx, w, http.StatusNotFound, ErrStringMetricNotFound)
		return
	}

	h.JSON(ctx, w, http.StatusOK, h.factory.BuildGetMetricResponse(ctx, metric))
}

func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextFromRequest(r)
	list, err := h.service.List(ctx)
	if err != nil {
		logger.New(ctx).Errorf("[metrics handler] error when getting metrics: %s", err.Error())
		h.PlainText(ctx, w, http.StatusInternalServerError, ErrStringDatabaseError)
		return
	}

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
