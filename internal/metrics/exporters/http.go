package exporters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	delivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type HTTPExporter struct {
	*worker.Worker
	address string
	factory delivery.MetricsRequestResponseFactory
	client  *resty.Client
}

func NewHTTPExporter(
	name string,
	factory delivery.MetricsRequestResponseFactory,
	cfg config.HTTPExporterConfig,
) *HTTPExporter {
	exp := &HTTPExporter{
		Worker:  worker.New(name, 1),
		address: cfg.Address,
		factory: factory,
		client: resty.New().
			SetTimeout(cfg.Timeout),
	}
	return exp
}

func (exp *HTTPExporter) Export(ctx context.Context, metrics []domain.Metric) error {
	req, err := exp.prepareRequest(ctx, metrics)
	if err != nil {
		return err
	}
	resp, err := req.Send()
	if err != nil {
		return err
	}
	logger.New(ctx).Infof("[http exporter] exported %d metrics successfully, status %s", len(metrics), resp.Status())
	return nil
}

func (exp *HTTPExporter) prepareRequest(ctx context.Context, metrics []domain.Metric) (*resty.Request, error) {
	// http://<АДРЕС_СЕРВЕРА>/updates
	body, err := json.Marshal(exp.factory.BuildUpdateBatchMetricRequest(ctx, metrics))
	if err != nil {
		return nil, err
	}

	req := exp.client.R().SetContext(ctx)
	req.URL = fmt.Sprintf("http://%s/updates", exp.address)
	req.Method = http.MethodPost
	req.SetBody(body)
	req.SetHeader("Content-Type", "application/json")

	return req, nil
}
