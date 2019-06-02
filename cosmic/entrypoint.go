package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
)

func CreateAndRun(applicationFactory func() *Application) {
	log.Init(log.DebugLevel, log.DebugLevel)

	// Build application with given factory
	app := applicationFactory()

	log.InfofCore("Starting application %s", app.Name)
	app.run()
}
