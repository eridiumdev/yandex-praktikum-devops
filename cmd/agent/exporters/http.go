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

type HTTPExporter struct {
	*AbstractExporter
	host   string
	port   int
	client *http.Client
}

func NewHTTPExporter(name string, host string, port int, timeout time.Duration) *HTTPExporter {
	exp := &HTTPExporter{
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
	exp.readyUp()
	return exp
}

func (exp *HTTPExporter) Export(ctx context.Context, mtx []metrics.Metric) error {
	defer func() {
		exp.readyUp()
	}()

	for _, metric := range mtx {
		req, err := exp.prepareRequest(ctx, metric)
		if err != nil {
			return err
		}
		err = exp.doRequest(req, metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (exp *HTTPExporter) prepareRequest(ctx context.Context, metric metrics.Metric) (*http.Request, error) {
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

func (exp *HTTPExporter) doRequest(request *http.Request, metric metrics.Metric) error {
	resp, err := exp.client.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	logger.Infof("export %s: %s", metric.GetName(), resp.Status)
	return nil
}
