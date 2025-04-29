package logger

import "log"

type Logger interface {
	Error(msg string, err error)
	Info(msg string)
}

type StdLogger struct{}

func NewLogger() Logger {
	return &StdLogger{}
}

func (l StdLogger) Error(msg string, err error) {
	log.Printf("[ERROR] %s: %v", msg, err)
}

func (l StdLogger) Info(msg string) {
	log.Printf("[INFO] %s", msg)
}

func (l StdLogger) Debug(msg string) {
	log.Printf("[DEBUG] %s", msg)
}
