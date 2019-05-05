package event

import "fmt"

// Signals mouse movement
type MouseMovedEvent struct {
	X, Y float32
}

func (e *MouseMovedEvent) Type() Type {
	return MouseMoved
}

func (e *MouseMovedEvent) Category() Category {
	return CategoryInput | CategoryMouse
}

func (e *MouseMovedEvent) String() string {
	return fmt.Sprintf("MouseMovedEvent [x=%f, y=%f]", e.X, e.Y)
}

// Base mouse button event, used for pressed and released events
type mouseButtonEvent struct {
	Button int
}

func (e *mouseButtonEvent) Category() Category {
	return CategoryInput | CategoryMouse | CategoryMouseButton
}

func (e *mouseButtonEvent) string(action string) string {
	return fmt.Sprintf("MouseButton%sEvent [button=%d]", action, e.Button)
}

// Signals that a mouse button was pressed
type MouseButtonPressedEvent struct {
	mouseButtonEvent
}

func (e *MouseButtonPressedEvent) Type() Type {
	return MouseButtonPressed
}

func (e *MouseButtonPressedEvent) String() string {
	return e.string("Pressed")
}

// Signals that a mouse button was released
type MouseButtonReleasedEvent struct {
	mouseButtonEvent
}

func (e *MouseButtonReleasedEvent) Type() Type {
	return MouseButtonReleased
}

func (e *MouseButtonReleasedEvent) String() string {
	return e.string("Released")
}

// Signals that the mouse wheel was scrolled
type MouseScrolledEvent struct {
	OffsetX, OffsetY float32
}

func (e *MouseScrolledEvent) Type() Type {
	return MouseScrolled
}

func (e *MouseScrolledEvent) Category() Category {
	return CategoryInput | CategoryMouse
}

func (e *MouseScrolledEvent) String() string {
	return fmt.Sprintf("MouseScrolledEvent [OffsetX=%f, OffsetY=%f]", e.OffsetX, e.OffsetY)
}
