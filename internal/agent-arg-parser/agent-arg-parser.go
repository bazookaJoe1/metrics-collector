package agent_arg_parser

import (
	"flag"
	"fmt"
	npv "github.com/bazookajoe1/metrics-collector/internal/net-params-validator"
	"strings"
	"time"
)

type CLArgsParams struct {
	Host           string
	Port           string
	PollInterval   time.Duration
	ReportInterval time.Duration
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

func (a *CLArgsParams) GetPI() time.Duration {
	return a.PollInterval
}

func (a *CLArgsParams) GetRI() time.Duration {
	return a.ReportInterval
}

func ArgParse() *CLArgsParams {
	params := &CLArgsParams{
		Host:           "localhost",
		Port:           "8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}

	flag.Var(params, "a", "Server listen point in format: `Host:Port`")
	flag.DurationVar(&params.ReportInterval, "r", 10*time.Second, "Report interval in seconds")
	flag.DurationVar(&params.PollInterval, "p", 2*time.Second, "Collect interval in seconds")
	flag.Parse()

	return params
}
