package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"runtime"
)

func init() {
	// Make sure main() runs on the main thread. TODO find out why exactly this is necessary
	runtime.LockOSThread()
}

func CreateAndRun(applicationFactory func() *Application, logLevelApp log.Level, logLevelCore log.Level) {
	log.Init(logLevelApp, logLevelCore)

	log.DebugCore("Creating application with given factory")

	// Build application with given factory
	app := applicationFactory()

	log.DebugfCore("Starting application %s", app.Name)
	app.run()
}
