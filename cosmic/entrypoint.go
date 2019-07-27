package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"runtime"
	"sync"
)

func init() {
	// Make sure main() runs on the main thread. TODO find out why exactly this is necessary
	runtime.LockOSThread()
}

// Certain systems of cosmic assume that only one instance of Application is
// created and active while running (an instance of) the engine. Instances of
// Application are thus created as a singleton.
var App *Application
var once sync.Once

func CreateAndRun(applicationFactory func() *Application, logLevelApp log.Level, logLevelCore log.Level) {
	once.Do(func() {
		log.Init(logLevelApp, logLevelCore)

		log.DebugCore("Creating application with given factory")

		// Build application with given factory
		App = applicationFactory()

		log.DebugfCore("Starting application %s", App.Name)
		App.run()
	})
}
