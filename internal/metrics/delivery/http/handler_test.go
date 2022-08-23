package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/routing"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/backup"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/hash"
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
		html += fmt.Sprintf("%s : %s", m.Name, m.StringValue())
	}
	html += "</html>"
	return []byte(html), nil
}

func getDummyRenderer() *dummyRenderer {
	return &dummyRenderer{}
}

func getDummyHasher() MetricsHasher {
	return hash.NewHasher("s3cr3t-k3y")
}

func getDummyFactory() MetricsRequestResponseFactory {
	return NewRequestResponseFactory(getDummyHasher())
}

func getDummyRepo() service.MetricsRepository {
	r := repository.NewInMemRepo()
	ctx := context.Background()
	_ = r.Store(ctx, domain.NewCounter(domain.PollCount, 5))
	_ = r.Store(ctx, domain.NewGauge(domain.Alloc, 10.123))
	return r
}

func getDummyBackuper() service.MetricsBackuper {
	return &backup.Mock{}
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
	body   string
	want   Want
}

func init() {
	logger.InitZerolog(context.Background(), config.LoggerConfig{
		Level: logger.LevelCritical,
		Mode:  logger.ModeDevelopment,
	})
}

func runTests(t *testing.T, tt TestCase) {
	ctx := context.Background()
	router := routing.NewChiRouter()
	repo := getDummyRepo()
	backuper := getDummyBackuper()

	svc, _ := service.NewMetricsService(ctx, repo, backuper, config.BackupConfig{
		Interval:  0,
		DoRestore: false,
	})

	h := NewMetricsHandler(svc, getDummyRenderer(), getDummyFactory(), getDummyHasher())
	router.AddRoute(http.MethodGet, "/", h.List)
	router.AddRoute(http.MethodPost, "/value", h.Get)
	router.AddRoute(http.MethodPost, "/update", h.Update)

	s := httptest.NewServer(router.Mux)
	defer s.Close()

	buffer := bytes.Buffer{}
	buffer.Write([]byte(tt.body))

	req, err := http.NewRequest(tt.method, s.URL+tt.url, &buffer)
	require.NoError(t, err)

	resp, doErr := http.DefaultClient.Do(req)
	require.NoError(t, doErr)

	body, readErr := io.ReadAll(resp.Body)
	require.NoError(t, readErr)
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
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter","delta":5}`,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"PollCount","type":"counter","delta":10,"hash":"953f4bbe05b8225241fbe00f63f2d84ae756165a81b25f89d96436c64f374793"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "positive test: gauge",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"Alloc","type":"gauge","value":10.20}`,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"Alloc","type":"gauge","value":10.2,"hash":"380f160a870d57bafccffae49443dd8d227c806485c2b4b80bcfaf96f19447b3"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "positive test: counter with hash",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter","delta":5,"hash":"7148ff92910a879bba42647839901cdd4f9c68f952657e36ead4e894511d82af"}`,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"PollCount","type":"counter","delta":10,"hash":"953f4bbe05b8225241fbe00f63f2d84ae756165a81b25f89d96436c64f374793"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "positive test: gauge with hash",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"Alloc","type":"gauge","value":10.20,"hash":"380f160a870d57bafccffae49443dd8d227c806485c2b4b80bcfaf96f19447b3"}`,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"Alloc","type":"gauge","value":10.2,"hash":"380f160a870d57bafccffae49443dd8d227c806485c2b4b80bcfaf96f19447b3"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad hash",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter","delta":5,"hash":"-"}`,
			want: Want{
				code:        http.StatusBadRequest,
				response:    ErrStringInvalidHash,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad url",
			url:    "/update/123",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter","delta":5}`,
			want: Want{
				code:        http.StatusNotFound,
				response:    "404 page not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad counter value",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter","delta":"abcd"}`,
			want: Want{
				code:        http.StatusBadRequest,
				response:    ErrStringInvalidJSON,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad gauge value",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"Alloc","type":"counter","delta":"abcd"}`,
			want: Want{
				code:        http.StatusBadRequest,
				response:    ErrStringInvalidJSON,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad metric type",
			url:    "/update",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"unknown"}`,
			want: Want{
				code:        http.StatusNotImplemented,
				response:    ErrStringInvalidMetricType,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: invalid JSON",
			url:    "/update",
			method: http.MethodPost,
			body:   `{123}`,
			want: Want{
				code:        http.StatusBadRequest,
				response:    ErrStringInvalidJSON,
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
			url:    "/value",
			body:   `{"id":"PollCount","type":"counter"}`,
			method: http.MethodPost,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"PollCount","type":"counter","delta":5,"hash":"7148ff92910a879bba42647839901cdd4f9c68f952657e36ead4e894511d82af"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "positive test: gauge",
			url:    "/value",
			method: http.MethodPost,
			body:   `{"id":"Alloc","type":"gauge"}`,
			want: Want{
				code:        http.StatusOK,
				response:    `{"id":"Alloc","type":"gauge","value":10.123,"hash":"7e9e3da35d6b5e7bd5b2458f14fd54f566cfeb0e5b192cc220a08cf0b42f14a3"}`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:   "negative test: bad url",
			url:    "/value/123",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"counter"}`,
			want: Want{
				code:        http.StatusNotFound,
				response:    "404 page not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: wrong metric type",
			url:    "/value",
			method: http.MethodPost,
			body:   `{"id":"PollCount","type":"unknown"}`,
			want: Want{
				code:        http.StatusNotImplemented,
				response:    ErrStringInvalidMetricType,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test: metric not found",
			url:    "/value",
			method: http.MethodPost,
			body:   `{"id":"abcd","type":"counter"}`,
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
			body:   "",
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
