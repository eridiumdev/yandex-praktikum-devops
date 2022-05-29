package exporters

import (
	"bytes"
	"context"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"fmt"
	"net/http"
	"time"
)

type HttpExporter struct {
	*AbstractExporter
	host   string
	port   int
	client *http.Client
}

func NewHttpExporter(name string, host string, port int, timeout time.Duration) *HttpExporter {
	exp := &HttpExporter{
		AbstractExporter: &AbstractExporter{
			name:  name,
			ready: make(chan bool),
		},
		host: host,
		port: port,
		client: &http.Client{
			Timeout: timeout,
		},
	}
	exp.makeReady()
	return exp
}

func (exp *HttpExporter) Export(mtx []metrics.Metric) error {
	defer func() {
		exp.makeReady()
	}()

	ctx := context.Background()

	for _, metric := range mtx {
		req, err := exp.prepareRequest(ctx, metric)
		if err != nil {
			return err
		}
		resp, err := exp.client.Do(req)
		if err != nil {
			return err
		}
		logger.Infof("export %s: %s", metric.GetName(), resp.Status)
	}

	return nil
}

func (exp *HttpExporter) prepareRequest(ctx context.Context, metric metrics.Metric) (*http.Request, error) {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	url := fmt.Sprintf("http://%s:%d/update/%s/%s/%s",
		exp.host,
		exp.port,
		metric.GetType(),
		metric.GetName(),
		metric.GetStringValue())

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "text/plain")

	return request, nil
}
