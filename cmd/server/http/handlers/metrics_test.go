package handlers

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/rendering"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type DummyRenderingEngine struct{}

func (e *DummyRenderingEngine) Render(templatePath string, data any) ([]byte, error) {
	return []byte(fmt.Sprintf("<html>%v</html>", data)), nil
}

func getDummyEngine() rendering.Engine {
	return &DummyRenderingEngine{}
}

func getDummyStorage() storage.Storage {
	s := storage.NewInMemStorage()
	_ = s.StoreMetric(metrics.CounterMetric{
		AbstractMetric: metrics.AbstractMetric{
			Name: metrics.PollCount,
		},
		Value: 5,
	})
	_ = s.StoreMetric(metrics.GaugeMetric{
		AbstractMetric: metrics.AbstractMetric{
			Name: metrics.Alloc,
		},
		Value: 10.123,
	})
	return s
}

func TestUpdate(t *testing.T) {
	type Want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want Want
	}{
		{
			name: "positive test: counter",
			url:  "/update/counter/PollCount/5",
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test: gauge",
			url:  "/update/gauge/Alloc/10.20",
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad url",
			url:  "/update/123",
			want: Want{
				code:        http.StatusNotFound,
				response:    "[metrics handler] bad request url, use /update/<metricType>/<metricName>/<metricValue>",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad counter value",
			url:  "/update/counter/PollCount/10.20",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value '10.20': strconv.ParseInt: parsing \"10.20\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad gauge value",
			url:  "/update/gauge/Alloc/abcd",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value 'abcd': strconv.ParseFloat: parsing \"abcd\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad metric type",
			url:  "/update/unknown/testCounter/100",
			want: Want{
				code:        http.StatusNotImplemented,
				response:    "[metrics handler] bad metric type 'unknown', use one of: gauge, counter",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)

			w := httptest.NewRecorder()
			h := NewMetricsHandler(storage.NewInMemStorage(), getDummyEngine())

			h.Update(w, request)
			resp := w.Result()

			assert.Equal(t, tt.want.code, resp.StatusCode, "status code")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content-type")

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "error when reading response body")

			if tt.want.contentType == "application/json" {
				assert.JSONEq(t, tt.want.response, string(respBody), "response")
			} else {
				assert.Equal(t, tt.want.response, strings.TrimRight(string(respBody), "\n"), "response")
			}
		})
	}
}

func TestGet(t *testing.T) {
	type Want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want Want
	}{
		{
			name: "positive test: counter",
			url:  "/value/counter/PollCount",
			want: Want{
				code:        http.StatusOK,
				response:    "5",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test: gauge",
			url:  "/value/gauge/Alloc",
			want: Want{
				code:        http.StatusOK,
				response:    "10.123",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad url",
			url:  "/value/123",
			want: Want{
				code:        http.StatusNotFound,
				response:    "[metrics handler] bad request url, use /value/<metricType>/<metricName>",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: wrong metric type",
			url:  "/value/counter/Alloc",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] error when getting metric from storage: 'counter': incorrect metric type",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: metric not found",
			url:  "/value/counter/abcd",
			want: Want{
				code:        http.StatusNotFound,
				response:    "[metrics handler] error when getting metric from storage: 'abcd': metric not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)

			w := httptest.NewRecorder()
			h := NewMetricsHandler(getDummyStorage(), getDummyEngine())

			h.Get(w, request)
			resp := w.Result()

			assert.Equal(t, tt.want.code, resp.StatusCode, "status code")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content-type")

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "error when reading response body")

			if tt.want.contentType == "application/json" {
				assert.JSONEq(t, tt.want.response, string(respBody), "response")
			} else {
				assert.Equal(t, tt.want.response, strings.TrimRight(string(respBody), "\n"), "response")
			}
		})
	}
}

func TestList(t *testing.T) {
	type Want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want Want
	}{
		{
			name: "positive test",
			url:  "/",
			want: Want{
				code:        http.StatusOK,
				response:    "<html>[{{Alloc} 10.123} {{PollCount} 5}]</html>",
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)

			w := httptest.NewRecorder()
			h := NewMetricsHandler(getDummyStorage(), getDummyEngine())

			h.List(w, request)
			resp := w.Result()

			assert.Equal(t, tt.want.code, resp.StatusCode, "status code")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content-type")

			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "error when reading response body")

			assert.Equal(t, tt.want.response, strings.TrimRight(string(respBody), "\n"),"response")
		})
	}
}
