package logging

// LogObserver defines the interface needed to record kupak logs
type LogObserver interface {
	// Record accepts an incoming log and a signal, which shows the log level
	Record(log string, sig int)
}
