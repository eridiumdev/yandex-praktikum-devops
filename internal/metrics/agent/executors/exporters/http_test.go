package exporters

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestPrepareRequest(t *testing.T) {
	type Want struct {
		url         string
		method      string
		body        interface{}
		contentType string
	}
	tests := []struct {
		name   string
		metric domain.Metric
		want   Want
	}{
		{
			name:   "counter",
			metric: domain.NewCounter(domain.PollCount, 5),
			want: Want{
				url:         "http://localhost:80/update/counter/PollCount/5",
				method:      http.MethodPost,
				body:        nil,
				contentType: "text/plain",
			},
		},
		{
			name:   "gauge",
			metric: domain.NewGauge(domain.Alloc, 10.123),
			want: Want{
				url:         "http://localhost:80/update/gauge/Alloc/10.123",
				method:      http.MethodPost,
				body:        nil,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewHTTPExporter("http", HTTPExporterSettings{Host: "localhost", Port: 80, Timeout: 5})
			req := exp.prepareRequest(context.Background(), tt.metric)

			assert.Equal(t, tt.want.url, req.URL, "url")
			assert.Equal(t, tt.want.method, req.Method, "method")
			assert.Equal(t, tt.want.body, req.Body, "body")
			assert.Equal(t, tt.want.contentType, req.Header.Get("Content-Type"), "content-type")
		})
	}
}
