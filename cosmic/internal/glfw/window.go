package glfw

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/graphics"
	"github.com/lentus/cosmic-engine/cosmic/internal/vulkan"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
)

// glfwWindow provides a cross-platform window implementation using glfw.
type glfwWindow struct {
	context      graphics.Context
	nativeWindow *glfw.Window
	title        string
	vsync        bool

	eventCallback func(e event.Event)
}

func NewWindow(title string, width, height int, graphicsProps graphics.ContextProperties) *glfwWindow {
	window := &glfwWindow{
		title: title,
		vsync: true,
	}

	var err error
	if err = glfw.Init(); err != nil {
		log.PanicfCore("Failed to initialise GLFW - %s", err.Error())
	}

	//glfw.WindowHint(glfw.Resizable, glfw.False)
	if window.vsync {
		glfw.WindowHint(glfw.RefreshRate, glfw.True)
	} else {
		glfw.WindowHint(glfw.RefreshRate, glfw.False)
	}

	window.nativeWindow, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		log.PanicfCore("Failed to create GLFW window - %s", err.Error())
	}

	window.setCallbacks()
	window.context = vulkan.NewContext(window.nativeWindow)

	return window
}

func (w *glfwWindow) setCallbacks() {
	w.nativeWindow.SetCloseCallback(func(window *glfw.Window) {
		w.eventCallback(&event.WindowClose{})
	})

	w.nativeWindow.SetFramebufferSizeCallback(func(window *glfw.Window, width int, height int) {
		w.context.SignalFramebufferResized()
		w.eventCallback(&event.WindowResize{Width: width, Height: height})
	})

	w.nativeWindow.SetPosCallback(func(window *glfw.Window, xpos int, ypos int) {
		w.eventCallback(&event.WindowMoved{X: float32(xpos), Y: float32(ypos)})
	})

	w.nativeWindow.SetFocusCallback(func(window *glfw.Window, focused bool) {
		if focused {
			w.eventCallback(&event.WindowFocus{})
		} else {
			w.eventCallback(&event.WindowLostFocus{})
		}
	})

	w.nativeWindow.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var e event.Event

		switch action {
		case glfw.Press:
			e = &event.KeyPressed{Key: FromNativeKey[key]}
		case glfw.Release:
			e = &event.KeyReleased{Key: FromNativeKey[key]}
		default: // glfw.Repeat
			e = &event.KeyPressed{Key: FromNativeKey[key], RepeatCount: 1}
		}

		w.eventCallback(e)
	})

	w.nativeWindow.SetCharCallback(func(window *glfw.Window, char rune) {
		w.eventCallback(&event.KeyTyped{Char: char})
	})

	w.nativeWindow.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			w.eventCallback(&event.MouseButtonPressed{Button: FromNativeMouseButton[button]})
		default: // glfw.Release
			w.eventCallback(&event.MouseButtonReleased{Button: FromNativeMouseButton[button]})
		}
	})

	w.nativeWindow.SetScrollCallback(func(window *glfw.Window, xoff float64, yoff float64) {
		w.eventCallback(&event.MouseScrolled{OffsetX: float32(xoff), OffsetY: float32(yoff)})
	})

	w.nativeWindow.SetCursorPosCallback(func(window *glfw.Window, xpos float64, ypos float64) {
		w.eventCallback(&event.MouseMoved{X: float32(xpos), Y: float32(ypos)})
	})
}

func (w *glfwWindow) OnUpdate() {
	glfw.PollEvents()
	w.context.Render()
}

func (w *glfwWindow) GetWidth() int {
	width, _ := w.nativeWindow.GetSize()
	return width
}

func (w *glfwWindow) GetHeight() int {
	_, height := w.nativeWindow.GetSize()
	return height
}

func (w *glfwWindow) SetEventCallback(callback func(e event.Event)) {
	w.eventCallback = callback
}

func (w *glfwWindow) SetVSync(vsync bool) {
	w.vsync = vsync
}

func (w *glfwWindow) IsVSync() bool {
	return w.vsync
}

func (w *glfwWindow) GetNativeWindow() interface{} {
	return w.nativeWindow
}

func (w *glfwWindow) Terminate() {
	w.context.Terminate()

	log.DebugCore("Terminating GLFW window")
	glfw.Terminate()
}
