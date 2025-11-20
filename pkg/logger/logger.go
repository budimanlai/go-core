package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	l.infoLogger.Printf(msg, fields...)
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	l.errorLogger.Printf(msg, fields...)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	l.debugLogger.Printf(msg, fields...)
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	l.warnLogger.Printf(msg, fields...)
}
