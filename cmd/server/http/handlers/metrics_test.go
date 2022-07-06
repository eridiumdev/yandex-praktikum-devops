package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/routers"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/rendering"
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
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

type Want struct {
	code        int
	response    string
	contentType string
}

type TestCase struct {
	name   string
	url    string
	method string
	body   io.Reader
	want   Want
}

func runTests(t *testing.T, tt TestCase) {
	r := routers.NewChiRouter()
	_ = NewMetricsHandler(r, getDummyStorage(), getDummyEngine())
	s := httptest.NewServer(r.Mux)
	defer s.Close()

	req, err := http.NewRequest(tt.method, s.URL+tt.url, tt.body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, tt.want.code, resp.StatusCode, "status code")
	assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content-type")

	if tt.want.contentType == "application/json" {
		assert.JSONEq(t, tt.want.response, string(body), "response")
	} else {
		assert.Equal(t, tt.want.response, strings.TrimRight(string(body), "\n"), "response")
	}
}

func TestUpdate(t *testing.T) {
	tests := []TestCase{
		{
			name:   "positive test: counter",
			url:    "/update/counter/PollCount/5",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "positive test: gauge",
			url:    "/update/gauge/Alloc/10.20",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad url",
			url:    "/update/123",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusNotFound,
				response:    "404 page not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad counter value",
			url:    "/update/counter/PollCount/10.20",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value '10.20': strconv.ParseInt: parsing \"10.20\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad gauge value",
			url:    "/update/gauge/Alloc/abcd",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value 'abcd': strconv.ParseFloat: parsing \"abcd\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad metric type",
			url:    "/update/unknown/testCounter/100",
			method: http.MethodPost,
			body:   nil,
			want: Want{
				code:        http.StatusNotImplemented,
				response:    "[metrics handler] bad metric type 'unknown', use one of: gauge, counter",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTests(t, tt)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []TestCase{
		{
			name:   "positive test: counter",
			url:    "/value/counter/PollCount",
			body:   nil,
			method: http.MethodGet,
			want: Want{
				code:        http.StatusOK,
				response:    "5",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "positive test: gauge",
			url:    "/value/gauge/Alloc",
			method: http.MethodGet,
			body:   nil,
			want: Want{
				code:        http.StatusOK,
				response:    "10.123",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad url",
			url:    "/value/123",
			method: http.MethodGet,
			body:   nil,
			want: Want{
				code:        http.StatusNotFound,
				response:    "404 page not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: wrong metric type",
			url:    "/value/counter/Alloc",
			method: http.MethodGet,
			body:   nil,
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] error when getting metric from storage: 'counter': incorrect metric type",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: metric not found",
			url:    "/value/counter/abcd",
			method: http.MethodGet,
			body:   nil,
			want: Want{
				code:        http.StatusNotFound,
				response:    "[metrics handler] error when getting metric from storage: 'abcd': metric not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTests(t, tt)
		})
	}
}

func TestList(t *testing.T) {
	tests := []TestCase{
		{
			name:   "positive test",
			url:    "/",
			method: http.MethodGet,
			body:   nil,
			want: Want{
				code:        http.StatusOK,
				response:    "<html>[{{Alloc} 10.123} {{PollCount} 5}]</html>",
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTests(t, tt)
		})
	}
}
