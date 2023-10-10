package agentargparser

import (
	"flag"
	"fmt"
	npv "github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"strings"
)

type CLArgsParams struct {
	Host           string
	Port           string
	PollInterval   int
	ReportInterval int
}

func (a *CLArgsParams) String() string {
	return a.Host + ":" + a.Port
}

func (a *CLArgsParams) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in a form Host:Port")
	}

	err := npv.ValidateIP(hp[0])
	if err != nil {
		return err
	}

	err = npv.ValidatePort(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = hp[1]
	return nil

}

func (a *CLArgsParams) GetAddr() string {
	return a.Host
}

func (a *CLArgsParams) GetPort() string {
	return a.Port
}

func (a *CLArgsParams) GetPI() int {
	return a.PollInterval
}

func (a *CLArgsParams) GetRI() int {
	return a.ReportInterval
}

func ArgParse() *CLArgsParams {
	params := &CLArgsParams{
		Host:           "localhost",
		Port:           "8080",
		ReportInterval: 10,
		PollInterval:   2,
	}

	flag.Var(params, "a", "Server listen point in format: `Host:Port`")
	flag.IntVar(&params.ReportInterval, "r", 10, "Report interval in seconds")
	flag.IntVar(&params.PollInterval, "p", 2, "Collect interval in seconds")
	flag.Parse()

	return params
}
