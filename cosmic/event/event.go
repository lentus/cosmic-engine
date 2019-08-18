package event

type Type int

const (
	// Application events
	TypeAppTick Type = iota
	TypeAppUpdate
	TypeAppRender

	// Window events
	TypeWindowClose
	TypeWindowResize
	TypeWindowFocus
	TypeWindowLostFocus
	TypeWindowMoved

	// Key events
	TypeKeyPressed
	TypeKeyReleased
	TypeKeyTyped

	// Mouse events
	TypeMouseButtonPressed
	TypeMouseButtonReleased
	TypeMouseMoved
	TypeMouseScrolled
)

type Category int

const (
	CategoryNone        Category = 0
	CategoryApplication          = 1 << 0
	CategoryWindow               = 1 << 1
	CategoryInput                = 1 << 2
	CategoryKey                  = 1 << 3
	CategoryMouse                = 1 << 4
	CategoryMouseButton          = 1 << 5
)

type Event interface {
	Type() Type // TODO redundant?
	Category() Category
	String() string

	IsHandled() bool
	SetHandled()
}

func IsInCategory(e Event, c Category) bool {
	return (e.Category() & c) != 0
}

type baseEvent struct {
	handled bool
}

func (e *baseEvent) IsHandled() bool {
	return e.handled
}

func (e *baseEvent) SetHandled() {
	e.handled = true
}
