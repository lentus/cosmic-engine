package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/input"
	"github.com/lentus/cosmic-engine/cosmic/internal/glfw"
	"github.com/lentus/cosmic-engine/cosmic/layer"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

type Application struct {
	Name        string
	WindowProps *WindowProperties

	layerStack layer.Stack
	window     window

	// Signals whether the application should close. Setting this to false
	// terminates the game loop next frame.
	running bool
}

func (app *Application) PushLayer(l layer.Layer) {
	app.layerStack.Push(l)
}

func (app *Application) PushOverlay(l layer.Layer) {
	app.layerStack.PushOverlay(l)
}

func (app *Application) run() {
	app.window = createWindow(app.WindowProps, app.onEvent)
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

func (app *Application) onEvent(e event.Event) {
	if !event.IsInCategory(e, event.CategoryInput) {
		log.DebugCore(e.String())
	}

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

func (app *Application) getNativeWindow() interface{} {
	return app.window.GetNativeWindow()
}

// Provides a way to query whether a key is being pressed without having to
// keep state in the application.
func IsKeyPressed(key input.Key) bool {
	return glfw.IsKeyPressed(key, App.getNativeWindow())
}

// Provides a way to query whether a mouse button is being pressed without
// having to keep state in the application.
func IsMouseButtonPressed(mouseButton input.MouseButton) bool {
	return glfw.IsMouseButtonPressed(mouseButton, App.getNativeWindow())
}
