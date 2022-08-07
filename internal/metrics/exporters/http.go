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

func (exp *HTTPExporter) Export(ctx context.Context, mtx []domain.Metric) error {
	for _, metric := range mtx {
		req, err := exp.prepareRequest(ctx, metric)
		if err != nil {
			return err
		}
		resp, err := req.Send()
		if err != nil {
			return err
		}
		logger.New(ctx).Infof("[http exporter] exported %s, status: %s", metric.Name, resp.Status())
	}
	return nil
}

func (exp *HTTPExporter) prepareRequest(ctx context.Context, metric domain.Metric) (*resty.Request, error) {
	// http://<АДРЕС_СЕРВЕРА>/update
	body, err := json.Marshal(exp.factory.BuildUpdateMetricRequest(ctx, metric))
	if err != nil {
		return nil, err
	}

	req := exp.client.R().SetContext(ctx)
	req.URL = fmt.Sprintf("http://%s/update", exp.address)
	req.Method = http.MethodPost
	req.SetBody(body)
	req.SetHeader("Content-Type", "application/json")

	return req, nil
}
