package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

type surface struct {
	ref vulkan.Surface

	capabilities vulkan.SurfaceCapabilities
	format       vulkan.SurfaceFormat
	presentMode  vulkan.PresentMode
}

func (ctx *Context) createSurface() {
	surfacePtr, err := ctx.nativeWindow.CreateWindowSurface(ctx.instance, nil)
	if err != nil {
		log.PanicfCore("failed to create vulkan window surface, %s", err.Error())
	}

	ctx.surface = surface{
		ref: vulkan.SurfaceFromPointer(surfacePtr),
	}
}

func (ctx *Context) initSurfaceProperties() {
	log.DebugCore("Creating window surface")

	ctx.surface.capabilities = getSurfaceCapabilities(ctx.gpu.ref, ctx.surface.ref)
	ctx.surface.capabilities.Deref()
	ctx.surface.capabilities.CurrentExtent.Deref()
	ctx.surface.capabilities.MinImageExtent.Deref()
	ctx.surface.capabilities.MaxImageExtent.Deref()

	surfaceFormats := getSurfaceFormats(ctx.gpu.ref, ctx.surface.ref)
	for i := range surfaceFormats {
		surfaceFormats[i].Deref()
	}
	ctx.surface.format = pickSurfaceFormat(surfaceFormats)

	supportedPresentModes := getPresentModes(ctx.gpu.ref, ctx.surface.ref)
	ctx.surface.presentMode = pickPresentMode(supportedPresentModes)
}

func getSurfaceCapabilities(gpu vulkan.PhysicalDevice, surface vulkan.Surface) vulkan.SurfaceCapabilities {
	surfaceCapabilities := vulkan.SurfaceCapabilities{}
	result := vulkan.GetPhysicalDeviceSurfaceCapabilities(gpu, surface, &surfaceCapabilities)
	panicOnError(result, "get surface capabilities")

	return surfaceCapabilities
}

func getSurfaceFormats(gpu vulkan.PhysicalDevice, surface vulkan.Surface) []vulkan.SurfaceFormat {
	var formatCount uint32
	result := vulkan.GetPhysicalDeviceSurfaceFormats(gpu, surface, &formatCount, nil)
	panicOnError(result, "get surface format count")
	if formatCount == 0 {
		log.PanicCore("no surface format found")
	}

	surfaceFormats := make([]vulkan.SurfaceFormat, formatCount)
	result = vulkan.GetPhysicalDeviceSurfaceFormats(gpu, surface, &formatCount, surfaceFormats)
	panicOnError(result, "get surface formats")

	return surfaceFormats
}

func getPresentModes(gpu vulkan.PhysicalDevice, surface vulkan.Surface) []vulkan.PresentMode {
	var presentModeCount uint32
	result := vulkan.GetPhysicalDeviceSurfacePresentModes(gpu, surface, &presentModeCount, nil)
	panicOnError(result, "retrieve supported present modes")
	supportedPresentModes := make([]vulkan.PresentMode, presentModeCount)
	result = vulkan.GetPhysicalDeviceSurfacePresentModes(gpu, surface, &presentModeCount, supportedPresentModes)
	panicOnError(result, "retrieve supported present modes")

	return supportedPresentModes
}

// pickSurfaceFormat attempts to use SRGB. If that is not supported, use the
// first format supported format.
func pickSurfaceFormat(supportedFormats []vulkan.SurfaceFormat) vulkan.SurfaceFormat {
	for _, format := range supportedFormats {
		if format.Format == vulkan.FormatB8g8r8a8Srgb && format.ColorSpace == vulkan.ColorSpaceSrgbNonlinear {
			return format
		}
	}

	return supportedFormats[0]
}

// pickPresentMode attempts to use Mailbox present mode if available. Otherwise,
// FIFO is used. THIS BEHAVIOUR ENABLES VSYNC BY DEFAULT! Use
// PresentModeImmediate to support disabled VSYNC.
func pickPresentMode(supportedPresentModes []vulkan.PresentMode) vulkan.PresentMode {
	for _, presentMode := range supportedPresentModes {
		if presentMode == vulkan.PresentModeMailbox {
			return presentMode
		}
	}

	return vulkan.PresentModeFifo
}
