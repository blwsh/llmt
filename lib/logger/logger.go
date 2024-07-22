package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
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
