package event

import "fmt"

// Provides common behaviour for key events
type keyEvent struct {
	baseEvent
}

func (e keyEvent) Category() Category {
	return CategoryInput & CategoryKey
}

func (e keyEvent) string(action string, keyCode int) string {
	return fmt.Sprintf("Key%sEvent [keycode=%d]", action, keyCode)
}

// Signals that a certain key was pressed
type KeyPressed struct {
	keyEvent

	KeyCode int
}

func (e KeyPressed) Type() Type {
	return TypeKeyPressed
}

func (e KeyPressed) String() string {
	return e.string("Pressed", e.KeyCode)
}

// Signals that a certain key was released
type KeyReleased struct {
	keyEvent

	KeyCode int
}

func (e KeyReleased) Type() Type {
	return TypeKeyReleased
}

func (e KeyReleased) String() string {
	return e.string("Released", e.KeyCode)
}

// Signals that a certain key was typed (pressed and released quickly)
type KeyTyped struct {
	keyEvent

	KeyCode int
}

func (e KeyTyped) Type() Type {
	return TypeKeyTyped
}

func (e KeyTyped) String() string {
	return e.string("Typed", e.KeyCode)
}
