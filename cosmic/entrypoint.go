package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
)

func Run(app *Application) {
	log.Init(log.DebugLevel, log.DebugLevel, app.Name)

	log.InfofCore("Starting application %s", app.Name)
	app.run()
}
