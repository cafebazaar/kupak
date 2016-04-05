package logging

import (
	GoLogging "github.com/op/go-logging"
)

type localLogger struct {
	logger *GoLogging.Logger
}

// Record is a part of LogObserver interface implementation
func (l *localLogger) Record(log string, sig int) {
	switch sig {
	case 0:
		l.logger.Debug(log)
	case 1:
		l.logger.Info(log)
	case 2:
		l.logger.Warning(log)
	case 3:
		l.logger.Error(log)
	}
}
