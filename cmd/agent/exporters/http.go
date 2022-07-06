package exporters

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type HTTPExporter struct {
	*AbstractExporter
	host   string
	port   int
	client *resty.Client
}

func NewHTTPExporter(name string, host string, port int, timeout time.Duration) *HTTPExporter {
	exp := &HTTPExporter{
		AbstractExporter: &AbstractExporter{
			name:  name,
			ready: make(chan bool),
		},
		host: host,
		port: port,
		client: resty.New().
			SetTimeout(timeout),
	}
	exp.readyUp()
	return exp
}

func (exp *HTTPExporter) Export(ctx context.Context, mtx []metrics.Metric) error {
	defer func() {
		exp.readyUp()
	}()

	for _, metric := range mtx {
		req := exp.prepareRequest(metric)
		resp, err := req.Send()
		logger.Infof("export %s: %s", metric.GetName(), resp.Status())
		if err != nil {
			return err
		}
	}
	return nil
}

func (exp *HTTPExporter) prepareRequest(metric metrics.Metric) *resty.Request {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	req := exp.client.R()
	req.URL = fmt.Sprintf("http://%s:%d/update/%s/%s/%s",
		exp.host,
		exp.port,
		metric.GetType(),
		metric.GetName(),
		metric.GetStringValue())
	req.Method = http.MethodPost
	req.SetHeader("Content-Type", "text/plain")

	return req
}
