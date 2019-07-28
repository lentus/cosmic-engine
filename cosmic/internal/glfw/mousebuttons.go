package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/lentus/cosmic-engine/cosmic/input"
)

var ToNativeMouseButton = map[input.MouseButton]glfw.MouseButton{
	input.MouseButton1: glfw.MouseButton1,
	input.MouseButton2: glfw.MouseButton2,
	input.MouseButton3: glfw.MouseButton3,
	input.MouseButton4: glfw.MouseButton4,
	input.MouseButton5: glfw.MouseButton5,
	input.MouseButton6: glfw.MouseButton6,
	input.MouseButton7: glfw.MouseButton7,
	input.MouseButton8: glfw.MouseButton8,
}

var FromNativeMouseButton = map[glfw.MouseButton]input.MouseButton{
	glfw.MouseButton1: input.MouseButton1,
	glfw.MouseButton2: input.MouseButton2,
	glfw.MouseButton3: input.MouseButton3,
	glfw.MouseButton4: input.MouseButton4,
	glfw.MouseButton5: input.MouseButton5,
	glfw.MouseButton6: input.MouseButton6,
	glfw.MouseButton7: input.MouseButton7,
	glfw.MouseButton8: input.MouseButton8,
}
