package logger

// Logger defines the contract for the logger than can be passed to the api client
type Logger interface {
	// Debug level logs
	Debugf(format string, args ...interface{})

	// Info level logs
	Infof(format string, args ...interface{})

	// Error level logs
	Errorf(format string, args ...interface{})

	// Fatal level logs
	Fatalf(format string, args ...interface{})
}
