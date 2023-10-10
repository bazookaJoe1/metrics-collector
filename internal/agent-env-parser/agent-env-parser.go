package agent_env_parser

import (
	npv "github.com/bazookajoe1/metrics-collector/internal/net-params-validator"
	"github.com/caarlos0/env/v9"
	"time"
)

const invalidDuration time.Duration = 123456789

type EnvParams struct {
	Address        []string `env:"ADDRESS" envSeparator:":"`
	Host           string
	Port           string
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func (a *EnvParams) GetAddr() string {
	return a.Host
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func (a *EnvParams) GetPI() time.Duration {
	return a.PollInterval
}

func (a *EnvParams) GetRI() time.Duration {
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
		err := npv.ValidateIP(a.Address[0])
		if err != nil {
			return
		}
		err = npv.ValidatePort(a.Address[1])
		if err != nil {
			return
		}
		a.Host = a.Address[0]
		a.Port = a.Address[1]
	}
}
