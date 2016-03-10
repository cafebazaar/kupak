package logging

// LogObserver defines the interface needed to record kupak logs
type LogObserver interface{    
    Record(log string)
}