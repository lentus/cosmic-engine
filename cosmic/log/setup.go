package log

import (
	"github.com/op/go-logging"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
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
	case LevelDebug:
		logging.SetLevel(logging.DEBUG, coreLogName)
	case LevelNotice:
		logging.SetLevel(logging.NOTICE, coreLogName)
	case LevelWarn:
		logging.SetLevel(logging.WARNING, coreLogName)
	case LevelError:
		logging.SetLevel(logging.ERROR, coreLogName)
	default:
		logging.SetLevel(logging.INFO, coreLogName)
	}

	// Init app logger
	appLog = logging.MustGetLogger(appLogName)
	switch appLevel {
	case LevelDebug:
		logging.SetLevel(logging.DEBUG, appLogName)
	case LevelNotice:
		logging.SetLevel(logging.NOTICE, appLogName)
	case LevelWarn:
		logging.SetLevel(logging.WARNING, appLogName)
	case LevelError:
		logging.SetLevel(logging.ERROR, appLogName)
	default:
		logging.SetLevel(logging.INFO, appLogName)
	}
}
