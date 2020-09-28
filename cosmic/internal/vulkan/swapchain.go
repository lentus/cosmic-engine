package vulkan

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

func (ctx *Context) createSwapchain() {
	log.DebugCore("Creating Vulkan swapchain")

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
	ctx.surface.presentMode = pickPresentMode(getPresentModes(ctx.gpu.ref, ctx.surface.ref))

	ctx.swapchainImageCount = determineImageCount(
		ctx.surface.capabilities.MinImageCount,
		ctx.surface.capabilities.MaxImageCount,
	)
	log.DebugfCore("Requesting %d swapchain images", ctx.swapchainImageCount)

	ctx.swapchainImageExtent = createImageExtent(ctx.surface.capabilities, ctx.nativeWindow)
	swapchainCreateInfo := vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          ctx.surface.ref,
		MinImageCount:    ctx.swapchainImageCount,
		ImageFormat:      ctx.surface.format.Format,
		ImageColorSpace:  ctx.surface.format.ColorSpace,
		ImageExtent:      ctx.swapchainImageExtent,
		ImageArrayLayers: 1, // No stereoscopic rendering, which requires 2
		ImageUsage:       vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		PreTransform:     ctx.surface.capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaOpaqueBit,
		PresentMode:      ctx.surface.presentMode,
		Clipped:          vulkan.True,
		OldSwapchain:     nil,
	}

	if ctx.gpu.queueFamilies.hasSeparatePresentQueue() {
		swapchainCreateInfo.ImageSharingMode = vulkan.SharingModeConcurrent
		swapchainCreateInfo.QueueFamilyIndexCount = 2
		swapchainCreateInfo.PQueueFamilyIndices = []uint32{
			ctx.gpu.queueFamilies.graphicsIndex,
			ctx.gpu.queueFamilies.presentIndex,
		}
	} else {
		swapchainCreateInfo.ImageSharingMode = vulkan.SharingModeExclusive
	}

	var swapchain vulkan.Swapchain
	result := vulkan.CreateSwapchain(ctx.device, &swapchainCreateInfo, nil, &swapchain)
	panicOnError(result, "create swapchain")
	ctx.swapchain = swapchain

	var swapchainImageCount uint32
	result = vulkan.GetSwapchainImages(ctx.device, ctx.swapchain, &swapchainImageCount, nil)
	panicOnError(result, "retrieve swapchain image count")
	ctx.swapchainImageCount = swapchainImageCount

	log.DebugfCore("Using %d swapchain images", ctx.swapchainImageCount)
}

func createImageExtent(capabilities vulkan.SurfaceCapabilities, nativeWindow *glfw.Window) vulkan.Extent2D {
	var swapchainImageExtent vulkan.Extent2D
	if capabilities.CurrentExtent.Width != vulkan.MaxUint32 {
		swapchainImageExtent.Width = capabilities.CurrentExtent.Width
		swapchainImageExtent.Height = capabilities.CurrentExtent.Height
	} else {
		width, height := nativeWindow.GetSize()
		swapchainImageExtent.Width = uint32(width)
		swapchainImageExtent.Height = uint32(height)

		if swapchainImageExtent.Width < capabilities.MinImageExtent.Width {
			swapchainImageExtent.Width = capabilities.MinImageExtent.Width
		}
		if swapchainImageExtent.Height < capabilities.MinImageExtent.Height {
			swapchainImageExtent.Height = capabilities.MinImageExtent.Height
		}
		if swapchainImageExtent.Width > capabilities.MaxImageExtent.Width {
			swapchainImageExtent.Width = capabilities.MaxImageExtent.Width
		}
		if swapchainImageExtent.Height > capabilities.MaxImageExtent.Height {
			swapchainImageExtent.Height = capabilities.MaxImageExtent.Height
		}
	}

	return swapchainImageExtent
}

func determineImageCount(min, max uint32) uint32 {
	// Try to use one more than the minimum (recommended)
	if min+1 > max {
		return max
	} else {
		return min + 1
	}
}

type imageResourceSet struct {
	image         vulkan.Image
	view          vulkan.ImageView
	commandBuffer vulkan.CommandBuffer
	framebuffer   vulkan.Framebuffer
}

func (ctx *Context) createSwapchainImages() {
	log.DebugCore("Creating Vulkan swapchain images")

	swapchainImages := make([]vulkan.Image, ctx.swapchainImageCount)
	result := vulkan.GetSwapchainImages(ctx.device, ctx.swapchain, &ctx.swapchainImageCount, swapchainImages)
	panicOnError(result, "create swapchain images")

	ctx.imageResourceSets = make([]imageResourceSet, ctx.swapchainImageCount)
	for i := range ctx.imageResourceSets {
		ctx.imageResourceSets[i].image = swapchainImages[i]

		imageViewCreateInfo := vulkan.ImageViewCreateInfo{
			SType:      vulkan.StructureTypeImageViewCreateInfo,
			Image:      ctx.imageResourceSets[i].image,
			ViewType:   vulkan.ImageViewType2d,
			Format:     ctx.surface.format.Format,
			Components: vulkan.ComponentMapping{}, // Use identity mapping for rgba components
			SubresourceRange: vulkan.ImageSubresourceRange{
				AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}

		var imageView vulkan.ImageView
		result = vulkan.CreateImageView(ctx.device, &imageViewCreateInfo, nil, &imageView)
		panicOnError(result, fmt.Sprintf("create swapchain image view nr %d", i))
		ctx.imageResourceSets[i].view = imageView
	}
}

func (ctx *Context) destroySwapchainImageViews() {
	for _, imageResourceSet := range ctx.imageResourceSets {
		vulkan.DestroyImageView(ctx.device, imageResourceSet.view, nil)
	}
}
