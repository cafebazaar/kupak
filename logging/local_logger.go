package logging

import (
	GoLogging "github.com/op/go-logging"
)

type localLogger struct {
	messageQueue chan struct {
		string
		int
	}
	logger *GoLogging.Logger
}

// Record is a part of LogObserver interface implementation
func (l *localLogger) Record(log string, sig int) {
	l.messageQueue <- struct {
		string
		int
	}{log, sig}
}

func (l *localLogger) start() {
	for v := range l.messageQueue {
		switch v.int {
		case 0:
			l.logger.Debug(v.string)
		case 1:
			l.logger.Info(v.string)
		case 2:
			l.logger.Warning(v.string)
		case 3:
			l.logger.Error(v.string)
		}
	}
}
