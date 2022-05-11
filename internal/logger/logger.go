package logger

import (
	"log"
	"os"
)

type LoggerLevel int

const (
	DEBUG LoggerLevel = iota
	INFO  LoggerLevel = iota
	ERROR LoggerLevel = iota
	FATAL LoggerLevel = iota
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

func New(level LoggerLevel) Logger {
	return &logger{
		level: level,
		debug: newLogger("DEBUG: "),
		info:  newLogger("INFO: "),
		error: newLogger("ERROR: "),
		fatal: newLogger("FATAL: "),
	}
}

type logger struct {
	level LoggerLevel
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

func (l *logger) Debug(msg string, args ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Printf(msg, args...)
	}
}

func (l *logger) Info(msg string, args ...interface{}) {
	if l.level <= INFO {
		l.info.Printf(msg, args...)
	}
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.error.Printf(msg, args...)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.fatal.Fatalf(msg, args...)
}

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.Ldate|log.Ltime)
}
