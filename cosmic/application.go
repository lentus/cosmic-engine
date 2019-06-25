package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/layer"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

type Application struct {
	Name        string
	WindowProps *WindowProperties

	layerStack layer.Stack
	window     Window

	// Signals whether the application should close. Setting this to false
	// terminates the game loop next frame.
	running bool
}

func (app *Application) run() {
	app.window = CreateWindow(app.WindowProps, app.OnEvent)
	defer app.window.Terminate()

	app.running = true
	for app.running {
		app.window.OnUpdate()

		// Update all layers
		for it := app.layerStack.Bottom(); it.Next(); {
			it.Get().OnUpdate()
		}
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

	// When getting a WindowClose event, signal the app to stop running.
	if _, ok := e.(event.WindowClose); ok {
		app.running = false
		return
	}

	// Otherwise, pass the event down the layerstack until it is handled.
	for it := app.layerStack.Top(); it.Prev(); {
		it.Get().OnEvent(e)

		if e.IsHandled() {
			break
		}
	}
}
