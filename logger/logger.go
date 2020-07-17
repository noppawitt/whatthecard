package logger

import "log"

// Logger is a logger
type Logger struct {
	level string
}

// NewLogger returns a new Logger
func NewLogger(level string) *Logger {
	return &Logger{
		level: level,
	}
}

// Debugf logs the message in the std printf format if the log level is set to debug
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level != "debug" {
		return
	}
	log.Printf(format, v...)
}

// Debug logs the message if the log level is set to debug
func (l *Logger) Debug(v ...interface{}) {
	if l.level != "debug" {
		return
	}
	log.Println(v...)
}

// Errorf always log the message in the std printf format
func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Errorf always log the message
func (l *Logger) Error(v ...interface{}) {
	log.Println(v...)
}
