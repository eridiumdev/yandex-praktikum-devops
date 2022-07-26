package exporters

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
				url:         "http://localhost:80/update",
				method:      http.MethodPost,
				body:        `{"id":"PollCount","type":"counter","delta":5}`,
				contentType: "application/json",
			},
		},
		{
			name:   "gauge",
			metric: domain.NewGauge(domain.Alloc, 10.333),
			want: Want{
				url:         "http://localhost:80/update",
				method:      http.MethodPost,
				body:        `{"id":"Alloc","type":"gauge","value":10.333}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewHTTPExporter("http", HTTPExporterSettings{Host: "localhost", Port: 80, Timeout: 5})
			req, err := exp.prepareRequest(context.Background(), tt.metric)
			require.NoError(t, err)

			assert.Equal(t, tt.want.url, req.URL, "url")
			assert.Equal(t, tt.want.method, req.Method, "method")
			assert.Equal(t, tt.want.body, string(req.Body.([]byte)), "body")
			assert.Equal(t, tt.want.contentType, req.Header.Get("Content-Type"), "content-type")
		})
	}
}
