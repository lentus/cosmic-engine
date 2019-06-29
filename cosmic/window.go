package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
)

type WindowApi string

type Window interface {
	OnUpdate()
	Terminate()

	GetWidth() int
	GetHeight() int
	IsVSync() bool
	SetVSync(vsync bool)

	setEventCallback(func(e event.Event))
}

type WindowProperties struct {
	Title  string
	Width  int
	Height int
	Api    WindowApi
}
