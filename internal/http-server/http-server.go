package HTTPServer

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/bazookajoe1/metrics-collector/internal/storage"
)

type HTTPServer struct {
	address string
	port    string
	router  *http.ServeMux
	strg    storage.Storage
}

func (serv *HTTPServer) Init(address string, port string, strg storage.Storage) error {
	err := isValidIP(address)
	if err != nil {
		return err
	}

	err = isValidPort(port)
	if err != nil {
		return err
	}

	serv.address, serv.port = address, port

	serv.router = http.NewServeMux()

	serv.strg = strg

	return nil
}

// Harness for router to register handlers
func (serv *HTTPServer) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	serv.router.HandleFunc(pattern, handler)
}

func (serv *HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", serv.address, serv.port)
	err := http.ListenAndServe(aP, serv.router)
	if err != nil {
		os.Exit(-1)
	}
}

func isValidIP(address string) error {
	valid := net.ParseIP(address)
	if (valid == nil) && (address != "localhost") {
		return errors.New("IP address is not valid")
	}

	return nil
}

func isValidPort(port string) error {
	uintPort, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return errors.New("invalid port number")
	}

	if uintPort > 65535 {
		return errors.New("invalid port number")
	}

	return nil
}
