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

func (serv *HttpServer) MetricHandler(res http.ResponseWriter, req *http.Request) {
	var err error

	if req.Method != http.MethodPost {
		http.Error(res, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pP := parsePath(req)
	/*if pP.metricType != "gauge" && pP.metricType != "counter" {
		http.Error(res, "", http.StatusBadRequest)
		return
	}*/

	if pP.metricValue == "" {
		http.Error(res, "404 Not Found", http.StatusNotFound)
		return
	}

	switch pP.metricType {
	case "gauge":
		err = serv.strg.UpdateGauge(pP.metricName, pP.metricValue)
		if err != nil {
			http.Error(res, "400 Bad Request", http.StatusBadRequest)
			return
		}

		res.Write([]byte{})
		return
	case "counter":
		err = serv.strg.UpdateCounter(pP.metricName, pP.metricValue)
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
