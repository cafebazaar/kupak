package logging

import (
	GoLogging "github.com/op/go-logging"
	"os"
)

func init() {
	goLoggingLogger := GoLogging.MustGetLogger("Kupak")
	localBackend := GoLogging.NewLogBackend(os.Stderr, "", 0)
	format := GoLogging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.3s} ▶ %{message} %{color:reset} `)
	localBackendFormatter := GoLogging.NewBackendFormatter(localBackend, format)
	GoLogging.SetBackend(localBackendFormatter)
	l := &localLogger{
		messageQueue: make(chan struct {
			string
			int
		}),
		logger: goLoggingLogger}
	go l.start()
	RegisterLogObserver(l)
}
