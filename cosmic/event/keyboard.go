package event

import "fmt"

// Provides common behaviour for key events
type keyEvent struct {
	baseEvent
}

func (e keyEvent) Category() Category {
	return CategoryInput & CategoryKey
}

// Signals that a certain key was pressed
type KeyPressed struct {
	keyEvent

	KeyCode     int
	RepeatCount int
}

func (e KeyPressed) Type() Type {
	return TypeKeyPressed
}

func (e KeyPressed) String() string {
	return fmt.Sprintf("KeyPressedEvent [keycode=%d, repeatCount=%d]", e.KeyCode, e.RepeatCount)
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
	return fmt.Sprintf("KeyReleasedEvent [keycode=%d]", e.KeyCode)
}

// Signals that a certain key was typed (pressed and released quickly)
type KeyTyped struct {
	keyEvent

	Char rune
}

func (e KeyTyped) Type() Type {
	return TypeKeyTyped
}

func (e KeyTyped) String() string {
	return fmt.Sprintf("KeyTypedEvent [keycode=%c]", e.Char)
}
