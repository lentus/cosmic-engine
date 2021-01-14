package graphics

type ContextProperties struct {
}

type Context interface {
	Render()
	Terminate()

	SignalFramebufferResized()
}
