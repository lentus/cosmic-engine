package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/layer"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

type Application struct {
	Name       string
	layerStack layer.Stack

	// Signals whether the application should close. Setting this to false
	// terminates the game loop next frame.
	running bool
}

func (app *Application) run() {
	app.running = true

	for app.running {
		// Update all layers
		app.layerStack.ForEachAscending(func(l layer.Layer) {
			l.OnUpdate()
		})

		log.DebugCore(app.Name)
		log.InfoCore(app.Name)
		log.NoticeCore(app.Name)
		log.WarnCore(app.Name)
		log.ErrorCore(app.Name)

		log.Debug(app.Name)
		log.Info(app.Name)
		log.Notice(app.Name)
		log.Warn(app.Name)
		log.Error(app.Name)

		app.running = false
	}
}

func (app *Application) PushLayer(l layer.Layer) {
	app.layerStack.Push(l)
}

func (app *Application) PushOverlay(l layer.Layer) {
	app.layerStack.PushOverlay(l)
}
