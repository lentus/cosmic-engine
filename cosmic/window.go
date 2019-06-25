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

	window.glfwWindow, err = glfw.CreateWindow(props.Width, props.Height, props.Title, nil, nil)
	if err != nil {
		glfw.Terminate()
		log.Panicf("Failed to create GLFW window - %s", err.Error())
	}

	window.glfwWindow.MakeContextCurrent()
	window.glfwWindow.SetCloseCallback(func(w *glfw.Window) {
		window.eventCallback(event.WindowClose{})
	})

	return window
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
