package event

import "fmt"

// Signals that a window should close
type WindowClose struct {
	baseEvent
}

func (e WindowClose) Type() Type {
	return TypeWindowClose
}

func (e WindowClose) Category() Category {
	return CategoryWindow
}

func (e WindowClose) String() string {
	return "WindowCloseEvent"
}

// Signals that a window was resized
type WindowResize struct {
	baseEvent

	Width, Height int
}

func (e WindowResize) Type() Type {
	return TypeWindowResize
}

func (e WindowResize) Category() Category {
	return CategoryWindow
}

func (e WindowResize) String() string {
	return fmt.Sprintf("WindowResizeEvent [width=%d, height=%d]", e.Width, e.Height)
}

// Signals that a window gained focus
type WindowFocus struct {
	baseEvent
}

func (e WindowFocus) Type() Type {
	return TypeWindowFocus
}

func (e WindowFocus) Category() Category {
	return CategoryWindow
}

func (e WindowFocus) String() string {
	return "WindowFocusEvent"
}

// Signals that a window lost focus
type WindowLostFocus struct {
	baseEvent
}

func (e WindowLostFocus) Type() Type {
	return TypeWindowLostFocus
}

func (e WindowLostFocus) Category() Category {
	return CategoryWindow
}

func (e WindowLostFocus) String() string {
	return "WindowLostFocusEvent"
}

// Signals that a window lost focus
type WindowMoved struct {
	baseEvent
}

func (e WindowMoved) Type() Type {
	return TypeWindowMoved
}

func (e WindowMoved) Category() Category {
	return CategoryWindow
}

func (e WindowMoved) String() string {
	return "WindowMovedEvent"
}
