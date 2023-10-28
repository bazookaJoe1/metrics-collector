package collector

// ILogger is the interfaces that allows to work with different loggers.
type ILogger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
}
