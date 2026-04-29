package logger

import (
	zap "go.uber.org/zap"
)

type MyLogger struct {
	Log *zap.SugaredLogger
}


func NewLogger() *MyLogger {

	log := zap.Must(zap.NewDevelopment())

	logger :=  log.Sugar()

	return &MyLogger{Log: logger}
}
