package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

const (
	WindowApiGlfw WindowApi = iota
)

func CreateWindow(props *WindowProperties, eventCallback func(e event.Event)) (window Window) {
	switch props.Api {
	case WindowApiGlfw:
		window = newGlfwWindow(props)
	default:
		log.Panicf("Invalid window API value %d, make sure this API is available on your platform", props.Api)
	}

	window.SetEventCallback(eventCallback)

	return
}
