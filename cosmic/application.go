package cosmic

type Application struct {
	Name string

	// Signals whether the application should close. Setting this to false
	// terminates the game loop next frame.
	shouldClose bool
}

func (app *Application) run() {
	for !app.shouldClose {

	}
}
