package exporters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/executor"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type HTTPExporter struct {
	*executor.Executor
	host   string
	port   int
	client *resty.Client
}

type HTTPExporterSettings struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func NewHTTPExporter(name string, settings HTTPExporterSettings) *HTTPExporter {
	exp := &HTTPExporter{
		Executor: executor.New(name),
		host:     settings.Host,
		port:     settings.Port,
		client: resty.New().
			SetTimeout(settings.Timeout),
	}
	exp.ReadyUp()
	return exp
}

func (exp *HTTPExporter) Export(ctx context.Context, mtx []domain.Metric) error {
	defer func() {
		exp.ReadyUp()
	}()

	for _, metric := range mtx {
		req, err := exp.prepareRequest(ctx, metric)
		if err != nil {
			return err
		}
		resp, err := req.Send()
		logger.New(ctx).Infof("[http exporter] exported %s, status: %s", metric.Name, resp.Status())
		if err != nil {
			return err
		}
	}
	return nil
}

func (exp *HTTPExporter) prepareRequest(ctx context.Context, metric domain.Metric) (*resty.Request, error) {
	// http://<АДРЕС_СЕРВЕРА>/update
	body, err := json.Marshal(domain.PrepareUpdateMetricRequest(metric))
	if err != nil {
		return nil, err
	}

	req := exp.client.R().SetContext(ctx)
	req.URL = fmt.Sprintf("http://%s:%d/update", exp.host, exp.port)
	req.Method = http.MethodPost
	req.SetBody(body)
	req.SetHeader("Content-Type", "application/json")

	return req, nil
}
