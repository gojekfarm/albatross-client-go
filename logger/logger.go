package logger

import "log"

// DefaultLogger is the default logger for the client.
// The client users should provide their own logger implmentations
type DefaultLogger struct{}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[Debug] "+format, args)
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[Info] "+format, args)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[Error] "+format, args)
}

func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	log.Printf("[Fatal] "+format, args)
}
