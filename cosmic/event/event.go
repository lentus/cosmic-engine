package event

type Type int

const (
	// Application events
	AppTick Type = iota
	AppUpdate
	AppRender

	// Window events
	WindowClosed
	WindowResized
	WindowFocus
	WindowLostFocus
	WindowMoved

	// Key events
	KeyPressed
	KeyReleased
	KeyTyped

	// Mouse events
	MouseButtonPressed
	MouseButtonReleased
	MouseMoved
	MouseScrolled
)

type Category int

const (
	None                Category = 0
	CategoryApplication          = 1 << 0
	CategoryWindow               = 1 << 1
	CategoryInput                = 1 << 2
	CategoryKey                  = 1 << 3
	CategoryMouse                = 1 << 4
	CategoryMouseButton          = 1 << 5
)

type Event interface {
	Type() Type
	Category() Category
	String() string
}

func IsInCategory(e Event, c Category) bool {
	return (e.Category() & c) != 0
}
