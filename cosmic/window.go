package cosmic

import "github.com/lentus/cosmic-engine/cosmic/event"

type WindowApi int

const (
	WindowApiGlfw WindowApi = iota
	WindowApiDxgi           // TODO Not yet implemented
)

type Window interface {
	OnUpdate()

	GetWidth() uint
	GetHeight() uint

	SetEventCallback(func(e event.Event))

	SetVSync(vsync bool)
	IsVSync() bool
}

type WindowProperties struct {
	Title  string
	Width  uint
	Height uint
	Api    WindowApi
}
