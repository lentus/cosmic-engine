package log

import (
	log "github.com/sirupsen/logrus"
)

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

func Init(level Level) {
	switch level {
	case TraceLevel:
		log.SetLevel(log.TraceLevel)
	case DebugLevel:
		log.SetLevel(log.DebugLevel)
	case WarnLevel:
		log.SetLevel(log.WarnLevel)
	case ErrorLevel:
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func Trace(format string, args ...interface{}) {
	log.Tracef(format, args)
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Panic(format string, args ...interface{}) {
	log.Panicf(format, args)
}
