package http

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

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/service"
)

type dummyRenderer struct{}

func (e *dummyRenderer) RenderList(list []domain.Metric) ([]byte, error) {
	html := "<html>"
	for i, m := range list {
		if i > 0 {
			html += " | "
		}
		html += fmt.Sprintf("%s : %s", m.Name(), m.StringValue())
	}
	html += "</html>"
	return []byte(html), nil
}

func getDummyRenderer() *dummyRenderer {
	return &dummyRenderer{}
}

func getDummyRepo() service.MetricsRepository {
	r := repository.NewInMemRepo()
	_ = r.Store(domain.NewCounter(domain.PollCount, 5))
	_ = r.Store(domain.NewGauge(domain.Alloc, 10.123))
	return r
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
	logger.Init(logger.LevelCritical, logger.ModeDevelopment)
	r := routing.NewChiRouter()
	svc := service.NewMetricsService(getDummyRepo())
	_ = NewMetricsHandler(r, svc, getDummyRenderer())
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
				response:    ErrStringInvalidMetricValue,
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
				response:    ErrStringInvalidMetricValue,
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
				response:    ErrStringInvalidMetricType,
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
				code:        http.StatusNotFound,
				response:    ErrStringMetricNotFound,
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
				response:    ErrStringMetricNotFound,
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
				response:    "<html>Alloc : 10.123 | PollCount : 5</html>",
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
