package agentenvparser

import (
	"github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"github.com/caarlos0/env/v9"
)

const invalidDuration int = 123456789

type EnvParams struct {
	Address        []string `env:"ADDRESS" envSeparator:":"`
	Host           string
	Port           string
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

func (a *EnvParams) GetAddr() string {
	return a.Host
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func (a *EnvParams) GetPI() int {
	return a.PollInterval
}

func (a *EnvParams) GetRI() int {
	return a.ReportInterval
}

func EnvParse() *EnvParams {
	var ep = &EnvParams{
		PollInterval:   invalidDuration,
		ReportInterval: invalidDuration,
	}

	err := env.Parse(ep)
	if err != nil {
		return nil
	}
	ep.splitHP()

	return ep
}

func (a *EnvParams) splitHP() {
	if len(a.Address) >= 2 {
		err := netparamsvalidator.ValidateIP(a.Address[0])
		if err != nil {
			return
		}
		err = netparamsvalidator.ValidatePort(a.Address[1])
		if err != nil {
			return
		}
		a.Host = a.Address[0]
		a.Port = a.Address[1]
	}
}
