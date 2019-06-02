package log

func Debug(msg string) {
	appLog.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	appLog.Debugf(format, args)
}

func Info(msg string) {
	appLog.Info(msg)
}

func Infof(format string, args ...interface{}) {
	appLog.Infof(format, args)
}

func Notice(msg string) {
	appLog.Notice(msg)
}

func Noticef(format string, args ...interface{}) {
	appLog.Noticef(format, args)
}

func Warn(msg string) {
	appLog.Warning(msg)
}

func Warnf(format string, args ...interface{}) {
	appLog.Warningf(format, args)
}

func Error(msg string) {
	appLog.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	appLog.Errorf(format, args)
}

func Critical(msg string) {
	appLog.Critical(msg)
}

func Criticalf(format string, args ...interface{}) {
	appLog.Criticalf(format, args)
}

func Panic(msg string) {
	appLog.Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	appLog.Panicf(format, args)
}
