package log

func DebugCore(msg string) {
	coreLog.Debug(msg)
}

func DebugfCore(format string, args ...interface{}) {
	coreLog.Debugf(format, args)
}

func InfoCore(msg string) {
	coreLog.Info(msg)
}

func InfofCore(format string, args ...interface{}) {
	coreLog.Infof(format, args)
}

func NoticeCore(msg string) {
	coreLog.Notice(msg)
}

func NoticefCore(format string, args ...interface{}) {
	coreLog.Noticef(format, args)
}

func WarnCore(msg string) {
	coreLog.Warning(msg)
}

func WarnfCore(format string, args ...interface{}) {
	coreLog.Warningf(format, args)
}

func ErrorCore(msg string) {
	coreLog.Error(msg)
}

func ErrorfCore(format string, args ...interface{}) {
	coreLog.Errorf(format, args)
}

func CriticalCore(msg string) {
	coreLog.Critical(msg)
}

func CriticalfCore(format string, args ...interface{}) {
	coreLog.Criticalf(format, args)
}

func PanicCore(msg string) {
	coreLog.Panic(msg)
}

func PanicfCore(format string, args ...interface{}) {
	coreLog.Panicf(format, args)
}
