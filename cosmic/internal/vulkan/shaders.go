package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/internal/vulkan/shaders"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
	"unsafe"
)

func (ctx Context) createShaderModule(shaderFileName string) vulkan.ShaderModule {
	shaderCode, err := shaders.Asset(shaderFileName)
	if err != nil {
		log.Panicf("failed to read shader file %s: %s", shaderFileName, err.Error())
	}

	shaderModuleCreateInfo := vulkan.ShaderModuleCreateInfo{
		SType:    vulkan.StructureTypeShaderModuleCreateInfo,
		Flags:    0,
		CodeSize: uint(len(shaderCode)),
		PCode:    sliceUint32(shaderCode),
	}
	var shaderModule vulkan.ShaderModule
	result := vulkan.CreateShaderModule(ctx.device, &shaderModuleCreateInfo, nil, &shaderModule)
	panicOnError(result, "create shader module")

	return shaderModule
}

// sliceUint32 and sliceHeader are copied from vulkan-go/asche, it's magic.
// https://github.com/vulkan-go/asche/blob/master/util.go#L179
func sliceUint32(data []byte) []uint32 {
	const m = 0x7fffffff
	return (*[m / 4]uint32)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&data)).Data))[:len(data)/4]
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
