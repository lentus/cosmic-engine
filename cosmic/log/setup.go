package log

import (
	"github.com/op/go-logging"
	"os"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	NoticeLevel
	WarnLevel
	ErrorLevel
)

var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{module:-4s} ▶ %{level:-7s}%{color:reset} %{message}",
)

const coreLogName = "CORE"
const appLogName = "APP"

var coreLog *logging.Logger
var appLog *logging.Logger

func Init(appLevel, coreLevel Level) {
	// Init stdErr logging backend
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, logFormat)
	logging.SetBackend(backendFormatter)

	// Init core logger
	coreLog = logging.MustGetLogger(coreLogName)
	switch coreLevel {
	case DebugLevel:
		logging.SetLevel(logging.DEBUG, coreLogName)
	case NoticeLevel:
		logging.SetLevel(logging.NOTICE, coreLogName)
	case WarnLevel:
		logging.SetLevel(logging.WARNING, coreLogName)
	case ErrorLevel:
		logging.SetLevel(logging.ERROR, coreLogName)
	default:
		logging.SetLevel(logging.INFO, coreLogName)
	}

	// Init app logger
	appLog = logging.MustGetLogger(appLogName)
	switch appLevel {
	case DebugLevel:
		logging.SetLevel(logging.DEBUG, appLogName)
	case NoticeLevel:
		logging.SetLevel(logging.NOTICE, appLogName)
	case WarnLevel:
		logging.SetLevel(logging.WARNING, appLogName)
	case ErrorLevel:
		logging.SetLevel(logging.ERROR, appLogName)
	default:
		logging.SetLevel(logging.INFO, appLogName)
	}
}
