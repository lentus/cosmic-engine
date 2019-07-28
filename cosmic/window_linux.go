package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/internal/glfw"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

const (
	WindowApiGlfw WindowApi = "GLFW"
)

func createWindow(props *WindowProperties, eventCallback func(e event.Event)) (window window) {
	log.DebugfCore("Creating %s window", props.Api)

	switch props.Api {
	case WindowApiGlfw:
		window = glfw.NewGlfwWindow(props.Title, props.Width, props.Height)
	default:
		log.PanicfCore("Invalid window API value %s, make sure this API is available on your platform", props.Api)
	}

	window.SetEventCallback(eventCallback)

	return
}
