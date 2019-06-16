package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/layer"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

type Application struct {
	Name string

	layerStack layer.Stack

	// Signals whether the application should close. Setting this to false
	// terminates the game loop next frame.
	running bool
}

func (app *Application) run() {
	app.running = true

	for app.running {
		// Update all layers
		for it := app.layerStack.Bottom(); it.Next(); {
			it.Get().OnUpdate()
		}

		app.running = false
	}
}

func (app *Application) PushLayer(l layer.Layer) {
	app.layerStack.Push(l)
}

func (app *Application) PushOverlay(l layer.Layer) {
	app.layerStack.PushOverlay(l)
}

func (app *Application) OnEvent(e event.Event) {
	log.DebugCore(e.String())

	for it := app.layerStack.Top(); it.Prev(); {
		it.Get().OnEvent(e)

		if e.IsHandled() {
			break
		}
	}
}
