package event

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/input"
)

// Provides common behaviour for key events
type keyEvent struct {
	baseEvent
}

func (e keyEvent) Category() Category {
	return CategoryInput | CategoryKey
}

// Signals that a certain key was pressed
type KeyPressed struct {
	keyEvent

	Key         input.Key
	RepeatCount int
}

func (e KeyPressed) Type() Type {
	return TypeKeyPressed
}

func (e KeyPressed) String() string {
	return fmt.Sprintf("KeyPressedEvent [keycode=%d, repeatCount=%d]", e.Key, e.RepeatCount)
}

// Signals that a certain key was released
type KeyReleased struct {
	keyEvent

	Key input.Key
}

func (e KeyReleased) Type() Type {
	return TypeKeyReleased
}

func (e KeyReleased) String() string {
	return fmt.Sprintf("KeyReleasedEvent [keycode=%d]", e.Key)
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
