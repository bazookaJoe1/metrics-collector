package logger

import (
	"go.uber.org/zap"
)

type ILogger interface {
	Info(string)
	Debug(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	z, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	z = z.WithOptions(zap.AddCallerSkip(1)) // we need this because of using harness to call log functions

	l := new(ZapLogger)
	l.logger = z

	return l
}

func (l *ZapLogger) Sync() {
	err := l.logger.Sync()
	if err != nil {
		panic(err)
	}
}

func (l *ZapLogger) Info(s string) {
	l.logger.Info(s)
}

func (l *ZapLogger) Debug(s string) {
	l.logger.Debug(s)
}

func (l *ZapLogger) Warn(s string) {
	l.logger.Warn(s)
}

func (l *ZapLogger) Error(s string) {
	l.logger.Error(s)
}

func (l *ZapLogger) Fatal(s string) {
	l.logger.Fatal(s)
}
