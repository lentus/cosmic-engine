package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

func (ctx *Context) initWindowSurface() {
	log.DebugCore("Creating window surface")

	surfaceCapabilities := vulkan.SurfaceCapabilities{}
	result := vulkan.GetPhysicalDeviceSurfaceCapabilities(ctx.gpu.ref, ctx.surface, &surfaceCapabilities)
	panicOnError(result, "get surface capabilities")

	surfaceCapabilities.Deref()
	surfaceCapabilities.CurrentExtent.Deref()
	surfaceCapabilities.MinImageExtent.Deref()
	surfaceCapabilities.MaxImageExtent.Deref()
	ctx.surfaceCapabilities = surfaceCapabilities

	var formatCount uint32
	result = vulkan.GetPhysicalDeviceSurfaceFormats(ctx.gpu.ref, ctx.surface, &formatCount, nil)
	panicOnError(result, "get physical device format count")
	if formatCount == 0 {
		log.PanicCore("no surface format found")
	}

	surfaceFormats := make([]vulkan.SurfaceFormat, formatCount)
	vulkan.GetPhysicalDeviceSurfaceFormats(ctx.gpu.ref, ctx.surface, &formatCount, surfaceFormats)
	for i := range surfaceFormats {
		surfaceFormats[i].Deref()
	}

	if surfaceFormats[0].Format == vulkan.FormatUndefined {
		ctx.surfaceFormat.Format = vulkan.FormatB8g8r8a8Unorm
		ctx.surfaceFormat.ColorSpace = vulkan.ColorSpaceSrgbNonlinear
	} else {
		ctx.surfaceFormat = surfaceFormats[0]
	}
}
