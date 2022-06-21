package handlers

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	type Want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		url string
		want   Want
	}{
		{
			name: "positive test: counter",
			url: "/update/counter/PollCount/5",
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test: gauge",
			url: "/update/gauge/Alloc/10.20",
			want: Want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad url",
			url: "/update/123",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad request url, use /update/<metricType>/<metricName>/<metricValue>",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad counter value",
			url: "/update/counter/PollCount/10.20",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value '10.20': strconv.ParseInt: parsing \"10.20\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: bad gauge value",
			url: "/update/gauge/Alloc/abcd",
			want: Want{
				code:        http.StatusBadRequest,
				response:    "[metrics handler] bad metric value 'abcd': strconv.ParseFloat: parsing \"abcd\": invalid syntax",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)

			w := httptest.NewRecorder()
			s := storage.NewInMemStorage()
			h := NewMetricsHandler(s)

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
