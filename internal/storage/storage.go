package storage

type Storage interface {
	Init()
	UpdateGauge(string, string) error
	UpdateCounter(string, string) error
}
