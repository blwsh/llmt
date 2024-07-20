package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) Logger
}

type logger struct {
	*zap.SugaredLogger
}

func New(isLiveMode bool) Logger {
	var (
		newLogger *zap.Logger
		err       error
	)

	if isLiveMode {
		newLogger, err = zap.NewProduction()
	} else {
		newLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	return logger{newLogger.Sugar()}
}

func (l logger) With(args ...interface{}) Logger {
	return logger{l.SugaredLogger.With(args...)}
}
