package cosmic

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/graphics"
)

type WindowApi string

type WindowProperties struct {
	Title  string
	Width  int
	Height int
	Api    WindowApi

	graphicsProperties graphics.ContextProperties
}

type window interface {
	OnUpdate()
	Terminate()

	GetWidth() int
	GetHeight() int
	IsVSync() bool
	SetVSync(vsync bool)
	GetNativeWindow() interface{}

	SetEventCallback(func(e event.Event))
}
