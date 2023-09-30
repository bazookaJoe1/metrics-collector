package httpserver

import (
	"net/http"
	"strings"
)

type parsedPath struct {
	hand        string
	metricType  string
	metricName  string
	metricValue string
}

func (serv *HTTPServer) MetricHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	serv.Logger.Println("Request", req.URL.Path)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if req.Method != http.MethodPost {
		http.Error(res, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pP := parsePath(req)

	if pP.metricValue == "" {
		http.Error(res, "404 Not Found", http.StatusNotFound)
		return
	}

	switch pP.metricType {
	case "gauge":
		err = serv.Strg.UpdateGauge(pP.metricName, pP.metricValue)
		if err != nil {
			http.Error(res, "400 Bad Request", http.StatusBadRequest)
			return
		}

		res.Write([]byte{})
		return
	case "counter":
		err = serv.Strg.UpdateCounter(pP.metricName, pP.metricValue)
		if err != nil {
			http.Error(res, "400 Bad Request", http.StatusBadRequest)
			return
		}

		res.Write([]byte{})
		return
	default:
		http.Error(res, "", http.StatusBadRequest)
	}
	_ = err
}

func parsePath(req *http.Request) parsedPath {
	path := req.URL.Path
	separatedPath := strings.Split(path, "/")
	separatedPath = separatedPath[1:]

	pP := parsedPath{
		hand: separatedPath[0],
	}

	switch len(separatedPath) {
	case 1:
		break
	case 2:
		pP.metricType = separatedPath[1]
	case 3:
		pP.metricType = separatedPath[1]
		pP.metricName = separatedPath[2]
	case 4:
		pP.metricType = separatedPath[1]
		pP.metricName = separatedPath[2]
		pP.metricValue = separatedPath[3]
	default:
		pP.metricType = separatedPath[1]
		pP.metricName = separatedPath[2]
		pP.metricValue = separatedPath[3]
	}

	return pP
}
