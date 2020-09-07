package logger

import (
	"fmt"
	"log"
)

// DefaultLogger is the default logger for the client.
// The client users should provide their own logger implmentations
type DefaultLogger struct{}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[Debug] %s", fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[Info] %s", fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[Error] %s", fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	log.Printf("[Fatal] %s", fmt.Sprintf(format, args...))
}
