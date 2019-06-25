package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

const (
	WindowApiGlfw WindowApi = iota
	WindowApiDxgi           // TODO Not yet implemented
)

func CreateWindow(props *WindowProperties, eventCallback func(e event.Event)) (window Window) {
	switch props.Api {
	case WindowApiGlfw:
		window = newGlfwWindow(props)
	case WindowApiDxgi:
		log.Panicf("DirectX windows are not implemented yet")
	default:
		log.Panicf("Invalid glfwWindow API value %d, make sure this API is available on your platform", props.Api)
	}

	window.SetEventCallback(eventCallback)

	return
}
