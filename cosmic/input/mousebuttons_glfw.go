package input

import "github.com/go-gl/glfw/v3.2/glfw"

var ToNativeMouseButton = map[MouseButton]glfw.MouseButton{
	MouseButton1: glfw.MouseButton1,
	MouseButton2: glfw.MouseButton2,
	MouseButton3: glfw.MouseButton3,
	MouseButton4: glfw.MouseButton4,
	MouseButton5: glfw.MouseButton5,
	MouseButton6: glfw.MouseButton6,
	MouseButton7: glfw.MouseButton7,
	MouseButton8: glfw.MouseButton8,
}

var FromNativeMouseButton = map[glfw.MouseButton]MouseButton{
	glfw.MouseButton1: MouseButton1,
	glfw.MouseButton2: MouseButton2,
	glfw.MouseButton3: MouseButton3,
	glfw.MouseButton4: MouseButton4,
	glfw.MouseButton5: MouseButton5,
	glfw.MouseButton6: MouseButton6,
	glfw.MouseButton7: MouseButton7,
	glfw.MouseButton8: MouseButton8,
}
