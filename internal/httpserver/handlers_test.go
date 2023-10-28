package httpserver

import (
	"bytes"
	"github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"
	"github.com/bazookajoe1/metrics-collector/internal/storage/memorystorage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHTTPServer. Some test can be passed only when all test cases sequentially executed.
func TestHTTPServer(t *testing.T) {
	storage := memorystorage.NewMemoryStorage(new(MockLogger), new(MockConfig))

	server := ServerNew(new(MockConfig), storage, new(MockLogger))

	server.InitRoutes()

	tests := []struct {
		name       string
		header     http.Header
		method     string
		path       string
		body       io.Reader
		expectCode int
		expectBody string
	}{
		{
			name:       "ok with saving gauge string",
			method:     http.MethodPost,
			path:       "/update/gauge/test/1",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: emptyString,
		},
		{
			name:       "ok with reading gauge string",
			method:     http.MethodGet,
			path:       "/value/gauge/test",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: "1",
		},
		{
			name:       "not found with reading non-existent gauge",
			method:     http.MethodGet,
			path:       "/value/gauge/non-existent",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "ok with saving counter string first",
			method:     http.MethodPost,
			path:       "/update/counter/test/40",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: emptyString,
		},
		{
			name:       "ok with reading counter string first",
			method:     http.MethodGet,
			path:       "/value/counter/test",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: "40",
		},
		{
			name:       "ok with saving counter string second",
			method:     http.MethodPost,
			path:       "/update/counter/test/40",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: emptyString,
		},
		{
			name:       "ok with reading counter string second",
			method:     http.MethodGet,
			path:       "/value/counter/test",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusOK,
			expectBody: "80",
		},
		{
			name:       "not found with reading non-existent counter",
			method:     http.MethodGet,
			path:       "/value/counter/non-existent",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "not found with saving gauge string",
			method:     http.MethodPost,
			path:       "/update/gauge/test",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "not found with saving gauge string",
			method:     http.MethodPost,
			path:       "/update/gauge/40",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "not found with saving counter string",
			method:     http.MethodPost,
			path:       "/update/counter/test",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "not found with saving counter string",
			method:     http.MethodPost,
			path:       "/update/counter/40",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusNotFound,
			expectBody: emptyString,
		},
		{
			name:       "bad request with saving gauge string",
			method:     http.MethodPost,
			path:       "/update/gauge/40/str",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusBadRequest,
			expectBody: emptyString,
		},
		{
			name:       "bad request with saving counter string",
			method:     http.MethodPost,
			path:       "/update/counter/40/str",
			body:       bytes.NewBuffer(nil),
			expectCode: http.StatusBadRequest,
			expectBody: emptyString,
		},
		{
			name:       "ok with saving gauge JSON",
			header:     http.Header{"Content-Type": {"application/json"}},
			method:     http.MethodPost,
			path:       "/update/",
			body:       bytes.NewBuffer([]byte("{\"id\":\"testJSON\",\"type\":\"gauge\",\"value\":10}")),
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"testJSON\",\"type\":\"gauge\",\"value\":10}",
		},
		{
			name:       "ok with reading gauge JSON",
			header:     http.Header{"Content-Type": {"application/json"}},
			method:     http.MethodPost,
			path:       "/value/",
			body:       bytes.NewBuffer([]byte("{\"id\":\"testJSON\",\"type\":\"gauge\"}")),
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"testJSON\",\"type\":\"gauge\",\"value\":10}",
		},
		{
			name:       "ok with saving counter JSON",
			header:     http.Header{"Content-Type": {"application/json"}},
			method:     http.MethodPost,
			path:       "/update/",
			body:       bytes.NewBuffer([]byte("{\"id\":\"test\",\"type\":\"counter\",\"delta\":40}")),
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"test\",\"type\":\"counter\",\"delta\":120}",
		},
		{
			name:       "ok with reading counter JSON",
			header:     http.Header{"Content-Type": {"application/json"}},
			method:     http.MethodPost,
			path:       "/value/",
			body:       bytes.NewBuffer([]byte("{\"id\":\"test\",\"type\":\"counter\"}")),
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"test\",\"type\":\"counter\",\"delta\":120}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body := request(tt.method, tt.path, tt.body, server.Server, tt.header)

			assert.Equal(t, tt.expectCode, code)
			assert.Equal(t, tt.expectBody, body)
		})
	}
}

type MockConfig struct{}

func (cm MockConfig) GetAddress() string                { return "" }
func (cm MockConfig) GetPort() string                   { return "" }
func (cm MockConfig) GetFileSaver() filesaver.FileSaver { return filesaver.FileSaver{} }

type MockLogger struct{}

func (l MockLogger) Info(string)  {}
func (l MockLogger) Error(string) {}
func (l MockLogger) Debug(string) {}
func (l MockLogger) Fatal(string) {}

func request(method, path string, body io.Reader, e *echo.Echo, header http.Header) (int, string) { // func taken from echo sources
	req := httptest.NewRequest(method, path, body)

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
