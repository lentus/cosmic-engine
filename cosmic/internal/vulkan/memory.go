package vulkan

import "github.com/vulkan-go/vulkan"

func (ctx *Context) findMemoryTypeIndex(requirements *vulkan.MemoryRequirements, properties vulkan.MemoryPropertyFlags) uint32 {
	for i := uint32(0); i < ctx.gpuMemoryProperties.MemoryTypeCount; i++ {
		// Find index of matching memory type
		if requirements.MemoryTypeBits&(1<<i) != 0 {
			ctx.gpuMemoryProperties.MemoryTypes[i].Deref()

			// Check that the matching memory type supports all required properties
			if ctx.gpuMemoryProperties.MemoryTypes[i].PropertyFlags&properties == properties {
				return i
			}
		}
	}

	// TODO Return error?
	return vulkan.MaxUint32
}
