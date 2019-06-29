package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

const (
	WindowApiGlfw WindowApi = "GLFW"
	WindowApiDxgi           = "DXGI" // TODO Not yet implemented
)

func CreateWindow(props *WindowProperties, eventCallback func(e event.Event)) (window Window) {
	log.DebugfCore("Creating %s window", props.Api)

	switch props.Api {
	case WindowApiGlfw:
		window = newGlfwWindow(props)
	case WindowApiDxgi:
		log.PanicfCore("DirectX window API is not implemented yet")
	default:
		log.PanicfCore("Invalid window API value %s, make sure this API is available on your platform", props.Api)
	}

	window.setEventCallback(eventCallback)

	return
}
