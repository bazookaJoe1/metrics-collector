package serverargparser

import (
	"flag"
	"fmt"
	"strings"

	npv "github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
)

type CLArgParams struct {
	Host string
	Port string
}

func (a *CLArgParams) String() string {
	return a.Host + ":" + a.Port
}

func (a *CLArgParams) Set(s string) error {
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

func (a *CLArgParams) GetAddr() string {
	return a.Host
}

func (a *CLArgParams) GetPort() string {
	return a.Port
}

func ArgParse() *CLArgParams {
	na := &CLArgParams{
		Host: "localhost",
		Port: "8080",
	}

	flag.Var(na, "a", "Server listen point in format: `Host:Port`")
	flag.Parse()

	return na
}
