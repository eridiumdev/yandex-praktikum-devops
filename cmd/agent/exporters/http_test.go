package exporters

import (
	"context"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	type Want struct {
		url         string
		method      string
		body        string
		contentType string
	}
	tests := []struct {
		name   string
		metric metrics.Metric
		want   Want
	}{
		{
			name: "counter",
			metric: metrics.CounterMetric{
				AbstractMetric: metrics.AbstractMetric{
					Name: metrics.PollCount,
				},
				Value: 5,
			},
			want: Want{
				url:         "http://localhost:80/update/counter/PollCount/5",
				method:      http.MethodPost,
				body:        "",
				contentType: "text/plain",
			},
		},
		{
			name: "gauge",
			metric: metrics.GaugeMetric{
				AbstractMetric: metrics.AbstractMetric{
					Name: metrics.Alloc,
				},
				Value: 10.123,
			},
			want: Want{
				url:         "http://localhost:80/update/gauge/Alloc/10.123000",
				method:      http.MethodPost,
				body:        "",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewHTTPExporter("http", "localhost", 80, 5)
			req, err := exp.prepareRequest(context.Background(), tt.metric)
			defer func() {
				if req != nil && req.Body != nil {
					_ = req.Body.Close()
				}
			}()
			body, _ := io.ReadAll(req.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.url, req.URL.String(), "url")
			assert.Equal(t, tt.want.method, req.Method, "method")
			assert.Equal(t, tt.want.body, string(body), "body")
			assert.Equal(t, tt.want.contentType, req.Header.Get("Content-Type"), "content-type")
		})
	}
}
