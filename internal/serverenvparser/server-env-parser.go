package serverenvparser

import (
	"github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"github.com/caarlos0/env/v9"
)

const invalidDuration int = 123456789

type EnvParams struct {
	Address       []string `env:"ADDRESS" envSeparator:":"`
	Host          string
	Port          string
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       *bool  `env:"RESTORE"`
}

func (a *EnvParams) GetAddr() string {
	return a.Host
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func (a *EnvParams) GetStoreInterval() int {
	return a.StoreInterval
}

func (a *EnvParams) GetFilePath() string {
	return a.FilePath
}

func (a *EnvParams) GetRestore() *bool {
	return a.Restore
}

func EnvParse() *EnvParams {
	var ep = &EnvParams{StoreInterval: invalidDuration}

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
