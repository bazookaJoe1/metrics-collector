package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	serverconfig "github.com/bazookajoe1/metrics-collector/internal/server-config"
	"github.com/bazookajoe1/metrics-collector/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, path string, method string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	logger := zlogger.NewZapLogger()
	servStorage := memstorage.NewInMemoryStorage()
	config := serverconfig.NewConfig(servStorage, logger)
	// TODO: init http server
	serv := ServerNew(config)
	serv.InitRoutes()

	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	var testTable = []struct {
		method string
		url    string
		want   string
		status int
	}{
		{http.MethodPost, "/update/gauge/somegauge/1.011", "", http.StatusOK},
		{http.MethodPost, "/update/counter/somecounter/1", "", http.StatusOK},
		{http.MethodPost, "/update/gauge/somegauge/str", "", http.StatusBadRequest},
		{http.MethodPost, "/update/counter/somecounter/str", "", http.StatusBadRequest},
		{http.MethodGet, "/update/counter/somecounter/str", "", http.StatusMethodNotAllowed},
		{http.MethodPost, "/update/counter/somecounter", "404 page not found\n", http.StatusNotFound},
	}
	for _, v := range testTable {
		resp, get := testRequest(t, ts, v.url, v.method)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
		resp.Body.Close()
	}
}
