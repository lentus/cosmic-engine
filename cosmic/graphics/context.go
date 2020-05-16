package graphics

type ImageBuffering int

const (
	DoubleBuffering ImageBuffering = iota
	TripleBuffering
)

type ContextProperties struct {
	BufferingType ImageBuffering
}

type Context interface {
	Render()
	Terminate()
}
