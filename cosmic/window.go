package cosmic

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/lentus/cosmic-engine/cosmic/event"
	"github.com/lentus/cosmic-engine/cosmic/log"
)

type WindowApi int

type Window interface {
	OnUpdate()

	GetWidth() int
	GetHeight() int

	SetEventCallback(func(e event.Event))

	SetVSync(vsync bool)
	IsVSync() bool

	Terminate()
}

type WindowProperties struct {
	Title  string
	Width  int
	Height int
	Api    WindowApi
}

// glfwWindow provides a cross-platform window implementation using glfw.
type glfwWindow struct {
	glfwWindow *glfw.Window
	title      string
	vsync      bool

	eventCallback func(e event.Event)
}

func newGlfwWindow(props *WindowProperties) *glfwWindow {
	window := &glfwWindow{
		title: props.Title,
		vsync: true,
	}

	var err error
	if err = glfw.Init(); err != nil {
		log.Panicf("Failed to initialise GLFW - %s", err.Error())
	}

	if window.vsync {
		glfw.WindowHint(glfw.RefreshRate, glfw.True)
	} else {
		glfw.WindowHint(glfw.RefreshRate, glfw.False)
	}

	window.glfwWindow, err = glfw.CreateWindow(props.Width, props.Height, props.Title, nil, nil)
	if err != nil {
		glfw.Terminate()
		log.Panicf("Failed to create GLFW window - %s", err.Error())
	}

	window.setCallbacks()
	window.glfwWindow.MakeContextCurrent()

	return window
}

func (w *glfwWindow) setCallbacks() {
	w.glfwWindow.SetCloseCallback(func(window *glfw.Window) {
		w.eventCallback(event.WindowClose{})
	})

	w.glfwWindow.SetSizeCallback(func(window *glfw.Window, width int, height int) {
		w.eventCallback(event.WindowResize{Width: width, Height: height})
	})

	w.glfwWindow.SetPosCallback(func(window *glfw.Window, xpos int, ypos int) {
		w.eventCallback(event.WindowMoved{X: float32(xpos), Y: float32(ypos)})
	})

	w.glfwWindow.SetFocusCallback(func(window *glfw.Window, focused bool) {
		if focused {
			w.eventCallback(event.WindowFocus{})
		} else {
			w.eventCallback(event.WindowLostFocus{})
		}
	})

	w.glfwWindow.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var e event.Event

		switch action {
		case glfw.Press:
			e = event.KeyPressed{KeyCode: int(key)}
		case glfw.Release:
			e = event.KeyReleased{KeyCode: int(key)}
		default: // glfw.Repeat
			e = event.KeyPressed{KeyCode: int(key), RepeatCount: 1}
		}

		w.eventCallback(e)
	})

	w.glfwWindow.SetCharCallback(func(window *glfw.Window, char rune) {
		w.eventCallback(event.KeyTyped{Char: char})
	})

	w.glfwWindow.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			w.eventCallback(event.MouseButtonPressed{Button: int(button)})
		default: // glfw.Release
			w.eventCallback(event.MouseButtonReleased{Button: int(button)})
		}
	})

	w.glfwWindow.SetScrollCallback(func(window *glfw.Window, xoff float64, yoff float64) {
		w.eventCallback(event.MouseScrolled{OffsetX: float32(xoff), OffsetY: float32(yoff)})
	})

	w.glfwWindow.SetCursorPosCallback(func(window *glfw.Window, xpos float64, ypos float64) {
		w.eventCallback(event.MouseMoved{X: float32(xpos), Y: float32(ypos)})
	})
}

func (w *glfwWindow) OnUpdate() {
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w *glfwWindow) GetWidth() int {
	width, _ := w.glfwWindow.GetSize()
	return width
}

func (w *glfwWindow) GetHeight() int {
	_, height := w.glfwWindow.GetSize()
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

func (w *glfwWindow) Terminate() {
	log.DebugCore("Terminating GLFW window")
	glfw.Terminate()
}
