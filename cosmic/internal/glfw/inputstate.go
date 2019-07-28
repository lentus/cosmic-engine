package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/lentus/cosmic-engine/cosmic/input"
)

// Provides a way to query whether a key is being pressed without having to
// keep state in the application.
//
// Accepts a native window to prevent circular imports with cosmic package.
func IsKeyPressed(key input.Key, nativeWindow interface{}) bool {
	return nativeWindow.(*glfw.Window).GetKey(ToNativeKey[key]) == glfw.Press
}

// Provides a way to query whether a mouse button is being pressed without
// having to keep state in the application.
//
// Accepts a native window to prevent circular imports with cosmic package.
func IsMouseButtonPressed(mouseButton input.MouseButton, nativeWindow interface{}) bool {
	return nativeWindow.(*glfw.Window).GetMouseButton(ToNativeMouseButton[mouseButton]) == glfw.Press
}
