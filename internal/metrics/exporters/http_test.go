package exporters

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/config"
	delivery "eridiumdev/yandex-praktikum-go-devops/internal/metrics/delivery/http"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/hash"
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
				body:        `{"id":"PollCount","type":"counter","delta":5,"hash":"7148ff92910a879bba42647839901cdd4f9c68f952657e36ead4e894511d82af"}`,
				contentType: "application/json",
			},
		},
		{
			name:   "gauge",
			metric: domain.NewGauge(domain.Alloc, 10.333),
			want: Want{
				url:         "http://localhost:80/update",
				method:      http.MethodPost,
				body:        `{"id":"Alloc","type":"gauge","value":10.333,"hash":"c4873e615e845fc90113575d072888a3f701c0620efb01bdce186d52ac1a3512"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewHTTPExporter("http",
				delivery.NewRequestResponseFactory(hash.NewHasher("s3cr3t-k3y")),
				config.HTTPExporterConfig{
					Address: "localhost:80",
					Timeout: 5,
				})
			req, err := exp.prepareRequest(context.Background(), tt.metric)
			require.NoError(t, err)

			assert.Equal(t, tt.want.url, req.URL, "url")
			assert.Equal(t, tt.want.method, req.Method, "method")
			assert.Equal(t, tt.want.body, string(req.Body.([]byte)), "body")
			assert.Equal(t, tt.want.contentType, req.Header.Get("Content-Type"), "content-type")
		})
	}
}
