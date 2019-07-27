package event

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/input"
)

// Signals mouse movement
type MouseMoved struct {
	baseEvent

	X, Y float32
}

func (e MouseMoved) Type() Type {
	return TypeMouseMoved
}

func (e MouseMoved) Category() Category {
	return CategoryInput | CategoryMouse
}

func (e MouseMoved) String() string {
	return fmt.Sprintf("MouseMovedEvent [x=%f, y=%f]", e.X, e.Y)
}

// Signals that the mouse wheel was scrolled
type MouseScrolled struct {
	baseEvent

	OffsetX, OffsetY float32
}

func (e MouseScrolled) Type() Type {
	return TypeMouseScrolled
}

func (e MouseScrolled) Category() Category {
	return CategoryInput | CategoryMouse
}

func (e MouseScrolled) String() string {
	return fmt.Sprintf("MouseScrolledEvent [OffsetX=%f, OffsetY=%f]", e.OffsetX, e.OffsetY)
}

// Provides common behaviour for mouse button events
type mouseButton struct {
	baseEvent
}

func (e mouseButton) Category() Category {
	return CategoryInput | CategoryMouse | CategoryMouseButton
}

func (e mouseButton) string(action string, button input.MouseButton) string {
	return fmt.Sprintf("MouseButton%sEvent [button=%d]", action, button)
}

// Signals that a mouse button was pressed
type MouseButtonPressed struct {
	mouseButton

	Button input.MouseButton
}

func (e MouseButtonPressed) Type() Type {
	return TypeMouseButtonPressed
}

func (e MouseButtonPressed) String() string {
	return e.string("Pressed", e.Button)
}

// Signals that a mouse button was released
type MouseButtonReleased struct {
	mouseButton

	Button input.MouseButton
}

func (e MouseButtonReleased) Type() Type {
	return TypeMouseButtonReleased
}

func (e MouseButtonReleased) String() string {
	return e.string("Released", e.Button)
}
