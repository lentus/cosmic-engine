package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
)

func Run(app *Application) {
	log.Init(log.DebugLevel)

	log.Info("Starting application %s", app.Name)
	app.run()
}
