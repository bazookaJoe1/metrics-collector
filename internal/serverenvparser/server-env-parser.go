package serverenvparser

import (
	npv "github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"github.com/caarlos0/env/v9"
)

type EnvParams struct {
	Address []string `env:"ADDRESS" envSeparator:":"`
	Host    string
	Port    string
}

func (a *EnvParams) GetAddr() string {
	return a.Host
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func EnvParse() *EnvParams {
	var ep = &EnvParams{}

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
