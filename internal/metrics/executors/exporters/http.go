package exporters

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/executor"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
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
		req := exp.prepareRequest(ctx, metric)
		resp, err := req.Send()
		logger.Infof("export %s: %s", metric.Name(), resp.Status())
		if err != nil {
			return err
		}
	}
	return nil
}

func (exp *HTTPExporter) prepareRequest(ctx context.Context, metric domain.Metric) *resty.Request {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	req := exp.client.R().SetContext(ctx)
	req.URL = fmt.Sprintf("http://%s:%d/update/%s/%s/%s",
		exp.host,
		exp.port,
		metric.Type(),
		metric.Name(),
		metric.StringValue())
	req.Method = http.MethodPost
	req.SetHeader("Content-Type", "text/plain")

	return req
}
