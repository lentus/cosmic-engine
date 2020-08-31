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
		ctx.swapchainImageCount, ctx.surface.capabilities.MinImageCount, ctx.surface.capabilities.MaxImageCount,
	)

	swapchainImageExtent := createImageExtent(ctx.surface.capabilities, ctx.nativeWindow)
	swapchainCreateInfo := vulkan.SwapchainCreateInfo{
		SType:                 vulkan.StructureTypeSwapchainCreateInfo,
		Surface:               ctx.surface.ref,
		MinImageCount:         ctx.swapchainImageCount,
		ImageFormat:           ctx.surface.format.Format,
		ImageColorSpace:       ctx.surface.format.ColorSpace,
		ImageExtent:           swapchainImageExtent,
		ImageArrayLayers:      1, // No stereoscopic rendering, which requires 2
		ImageUsage:            vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		ImageSharingMode:      vulkan.SharingModeExclusive,
		QueueFamilyIndexCount: 0,   // Ignored since sharing mode is exclusive
		PQueueFamilyIndices:   nil, // Ignored since sharing mode is exclusive
		PreTransform:          vulkan.SurfaceTransformIdentityBit,
		CompositeAlpha:        vulkan.CompositeAlphaOpaqueBit,
		PresentMode:           ctx.surface.presentMode,
		Clipped:               vulkan.True,
		OldSwapchain:          nil,
	}
	var swapchain vulkan.Swapchain
	result := vulkan.CreateSwapchain(ctx.device, &swapchainCreateInfo, nil, &swapchain)
	panicOnError(result, "create swapchain")

	ctx.swapchain = swapchain

	var swapchainImageCount uint32
	result = vulkan.GetSwapchainImages(ctx.device, ctx.swapchain, &swapchainImageCount, nil)
	panicOnError(result, "retrieve swapchain image count")

	ctx.swapchainImageCount = swapchainImageCount
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

func determineImageCount(requested, min, max uint32) uint32 {
	if requested < min+1 {
		log.WarnfCore("Requested image count %d not supported by your system, min is %d", requested, min+1)
		return min + 1
	}

	if max > 0 && requested > max {
		log.WarnfCore("Requested image count %d not supported by your system, max is %d", requested, max)
		return max
	}

	return requested
}

func (ctx *Context) createSwapchainImages() {
	log.DebugCore("Creating Vulkan swapchain images")
	ctx.swapchainImages = make([]vulkan.Image, ctx.swapchainImageCount)
	result := vulkan.GetSwapchainImages(ctx.device, ctx.swapchain, &ctx.swapchainImageCount, ctx.swapchainImages)
	panicOnError(result, "create swapchain images")

	ctx.swapchainImageViews = make([]vulkan.ImageView, ctx.swapchainImageCount)
	for i := range ctx.swapchainImageViews {
		imageViewCreateInfo := vulkan.ImageViewCreateInfo{
			SType:      vulkan.StructureTypeImageViewCreateInfo,
			Image:      ctx.swapchainImages[i],
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
		ctx.swapchainImageViews[i] = imageView
	}
}

func (ctx *Context) destroySwapchainImageViews() {
	for _, imageView := range ctx.swapchainImageViews {
		vulkan.DestroyImageView(ctx.device, imageView, nil)
	}
}
