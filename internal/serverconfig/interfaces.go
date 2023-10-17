package serverconfig

import "github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"

// IConfig is an interface providing methods to configure httpserver.HTTPServer params
type IConfig interface {
	GetAddress() string
	GetPort() string
	GetFileSaver() filesaver.FileSaver
}
