package serverargparser

import (
	"flag"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"os"
	"strings"
)

type CLArgParams struct {
	Host          string
	Port          string
	StoreInterval int
	FilePath      string
	Restore       *bool
}

func (a *CLArgParams) String() string {
	return a.Host + ":" + a.Port
}

func (a *CLArgParams) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in a form Host:Port")
	}

	err := netparamsvalidator.ValidateIP(hp[0])
	if err != nil {
		return err
	}

	err = netparamsvalidator.ValidatePort(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = hp[1]
	return nil

}

func (a *CLArgParams) GetAddr() string {
	return a.Host
}

func (a *CLArgParams) GetPort() string {
	return a.Port
}

func (a *CLArgParams) GetStoreInterval() int {
	return a.StoreInterval
}

func (a *CLArgParams) GetFilePath() string {
	return a.FilePath
}

func (a *CLArgParams) GetRestore() *bool {
	return a.Restore
}

func ArgParse() *CLArgParams {
	na := &CLArgParams{
		Host:          "localhost",
		Port:          "8080",
		StoreInterval: 300,
		FilePath:      os.TempDir() + "\\metrics-db.json",
		Restore:       new(bool),
	}

	flag.Var(na, "a", "Server listen point in format: `Host:Port`")
	flag.IntVar(&na.StoreInterval, "i", 300, "Save metrics to file interval")
	flag.StringVar(&na.FilePath, "f", os.TempDir()+"\\metrics-db.json", "File name where to store metrics")
	flag.BoolVar(na.Restore, "r", true, "Load storage from file or not")
	flag.Parse()

	return na
}
