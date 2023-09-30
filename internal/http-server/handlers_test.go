package httpserver

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bazookajoe1/metrics-collector/internal/storages/memstorage"
)

func TestHTTPServer_MetricHandler(t *testing.T) {
	logger := log.New(os.Stderr, "", 0)
	servStorage := &memstorage.InMemoryStorage{}
	servStorage.Init()
	serv := &HTTPServer{
		Address: "",
		Port:    "",
		Router:  http.NewServeMux(),
		Strg:    servStorage,
		Logger:  logger,
	}
	serv.Router.HandleFunc("/update/", serv.MetricHandler)

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "Test OK gauge",
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/somegauge/1.011", nil),
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Test OK counter",
			request: httptest.NewRequest(http.MethodPost, "/update/counter/somecounter/1", nil),
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Test FAIL gauge",
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/somegauge/str", nil),
			want: want{
				code:        400,
				response:    "400 Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Test FAIL counter",
			request: httptest.NewRequest(http.MethodPost, "/update/counter/somecounter/str", nil),
			want: want{
				code:        400,
				response:    "400 Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Test INCORRECT method",
			request: httptest.NewRequest(http.MethodGet, "/update/counter/somecounter/str", nil),
			want: want{
				code:        405,
				response:    "405 Method not allowed",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Test NO value",
			request: httptest.NewRequest(http.MethodPost, "/update/counter/somecounter", nil),
			want: want{
				code:        404,
				response:    "404 Not Found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			serv.MetricHandler(w, test.request)

			res := w.Result()

			if !reflect.DeepEqual(res.StatusCode, test.want.code) {
				t.Errorf("Status code got: %v, want: %v", res.StatusCode, test.want.code)
			}
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			if err != nil {
				t.Fatalf("Cannot read response body: %v", err)
			}

			if !reflect.DeepEqual(strings.TrimSuffix(string(resBody), "\n"), test.want.response) {
				t.Errorf("Body got: %v, want: %v", string(resBody), test.want.response)
			}

			if !reflect.DeepEqual(res.Header.Get("Content-Type"), test.want.contentType) {
				t.Errorf("Content type got: %v, want: %v", res.Header.Get("Content-Type"), test.want.contentType)
			}
		})
	}
}
