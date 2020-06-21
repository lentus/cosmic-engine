package vulkan

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/graphics"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
	"unsafe"
)

type Context struct {
	nativeWindow *glfw.Window

	instance      vulkan.Instance
	gpu           physicalDevice
	device        vulkan.Device
	graphicsQueue vulkan.Queue
	presentQueue  vulkan.Queue

	surface             vulkan.Surface
	surfaceFormat       vulkan.SurfaceFormat
	surfaceCapabilities vulkan.SurfaceCapabilities

	swapchain                 vulkan.Swapchain
	swapchainImageCount       uint32
	swapchainImages           []vulkan.Image
	swapchainImageViews       []vulkan.ImageView
	activeSwapchainImageindex uint32

	swapchainImageAvailable vulkan.Fence

	framebuffers []vulkan.Framebuffer
	renderPass   vulkan.RenderPass

	depthStencilFormat      vulkan.Format
	depthStencilImage       vulkan.Image
	depthStencilImageMemory vulkan.DeviceMemory
	depthStencilImageView   vulkan.ImageView
	stencilAvailable        bool

	availableInstanceLayers     []vulkan.LayerProperties
	availableInstanceExtensions []vulkan.ExtensionProperties
	availableDeviceExtensions   []vulkan.ExtensionProperties
	enabledInstanceLayers       []string
	enabledInstanceExtensions   []string
	enabledDeviceExtensions     []string

	debugCallback vulkan.DebugReportCallback

	// These are only here for demo purposes so I can test the vulkan implementation
	commandPool             vulkan.CommandPool
	commandBuffer           vulkan.CommandBuffer
	renderCompleteSemaphore vulkan.Semaphore
}

func NewContext(nativeWindow *glfw.Window, bufferingType graphics.ImageBuffering) *Context {
	log.InfoCore("Creating Vulkan graphics context")

	if !glfw.VulkanSupported() {
		log.PanicCore("glfw reports that Vulkan is not supported, aborting")
	}

	ctx := Context{
		nativeWindow:              nativeWindow,
		enabledInstanceLayers:     make([]string, 0),
		enabledInstanceExtensions: make([]string, 0),
		enabledDeviceExtensions:   make([]string, 0),
		swapchainImageCount:       uint32(bufferingType) + 2, // First buffering type DoubleBuffering has index 0
	}
	ctx.nativeWindow.MakeContextCurrent()

	vulkan.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vulkan.Init(); err != nil {
		log.PanicfCore("failed to initialise Vulkan, %s", err.Error())
	}

	ctx.setupLayersAndExtensions()
	ctx.setupDebug()
	ctx.createVulkanInstance()
	ctx.initDebugCallback()
	ctx.createSurface()
	ctx.selectPhysicalDevice()
	ctx.createLogicalDevice()
	ctx.initWindowSurface()
	ctx.createSwapchain()
	ctx.createSwapchainImages()
	ctx.createDepthStencilImage()
	ctx.createRenderPass()
	ctx.createFramebuffers()
	ctx.createSynchronizations()

	// These are only here so I can test the vulkan implementation
	ctx.createCommandPool()
	ctx.createCommandBuffer()
	ctx.createRenderCompleteSemaphore()

	return &ctx
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.QueueWaitIdle(ctx.graphicsQueue) // Wait for the graphics queue to be idle

	vulkan.DestroySemaphore(ctx.device, ctx.renderCompleteSemaphore, nil)
	vulkan.DestroyCommandPool(ctx.device, ctx.commandPool, nil)

	vulkan.DestroyFence(ctx.device, ctx.swapchainImageAvailable, nil)
	ctx.destroyFramebuffers()
	vulkan.DestroyRenderPass(ctx.device, ctx.renderPass, nil)
	ctx.destroyDepthStencilImage()
	ctx.destroySwapchainImageViews()
	vulkan.DestroySwapchain(ctx.device, ctx.swapchain, nil) // Destroys swapchain images as well
	vulkan.DestroySurface(ctx.instance, ctx.surface, nil)
	vulkan.DestroyDevice(ctx.device, nil)
	vulkan.DestroyDebugReportCallback(ctx.instance, ctx.debugCallback, nil)
	vulkan.DestroyInstance(ctx.instance, nil)
}

func (ctx *Context) createVulkanInstance() {
	log.DebugCore("Creating Vulkan instance")

	// TODO get version info and application name from somewhere
	applicationInfo := vulkan.ApplicationInfo{
		SType:         vulkan.StructureTypeApplicationInfo,
		ApiVersion:    vulkan.MakeVersion(1, 1, 88),
		PEngineName:   safeStr("Cosmic Engine"),
		EngineVersion: vulkan.MakeVersion(0, 1, 0),
	}

	debugReportCallbackCreateInfo := createDebugReportCallbackCreateInfo()
	instanceCreateInfo := vulkan.InstanceCreateInfo{
		SType:                   vulkan.StructureTypeInstanceCreateInfo,
		PNext:                   unsafe.Pointer(&debugReportCallbackCreateInfo), // Enable debug callback for instance creation
		PApplicationInfo:        &applicationInfo,
		EnabledLayerCount:       uint32(len(ctx.enabledInstanceLayers)),
		PpEnabledLayerNames:     ctx.enabledInstanceLayers,
		EnabledExtensionCount:   uint32(len(ctx.enabledInstanceExtensions)),
		PpEnabledExtensionNames: ctx.enabledInstanceExtensions,
	}

	var instance vulkan.Instance
	result := vulkan.CreateInstance(&instanceCreateInfo, nil, &instance)
	panicOnError(result, "create Vulkan instance")

	ctx.instance = instance
	if err := vulkan.InitInstance(ctx.instance); err != nil {
		log.PanicfCore("failed to initialise Vulkan instance (%s)", err.Error())
	}
}

func (ctx *Context) createSurface() {
	surfacePtr, err := ctx.nativeWindow.CreateWindowSurface(ctx.instance, nil)
	if err != nil {
		log.PanicfCore("failed to create vulkan window surface, %s", err.Error())
	}
	ctx.surface = vulkan.SurfaceFromPointer(surfacePtr)
}

func (ctx *Context) createLogicalDevice() {
	log.DebugCore("Creating Vulkan device")

	uniqueQueueFamilyIndices := []uint32{ctx.gpu.queueFamilies.graphicsIndex}
	if ctx.gpu.queueFamilies.hasSeparatePresentQueue() {
		uniqueQueueFamilyIndices = append(uniqueQueueFamilyIndices, ctx.gpu.queueFamilies.presentIndex)
	}

	var queueCreateInfos []vulkan.DeviceQueueCreateInfo
	for _, queueFamilyIndex := range uniqueQueueFamilyIndices {
		queueCreateInfos = append(queueCreateInfos, vulkan.DeviceQueueCreateInfo{
			SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: queueFamilyIndex,
			QueueCount:       1,
			PQueuePriorities: []float32{1.0},
		})
	}

	deviceCreateInfo := vulkan.DeviceCreateInfo{
		SType:                   vulkan.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    uint32(len(queueCreateInfos)),
		PQueueCreateInfos:       queueCreateInfos,
		EnabledExtensionCount:   uint32(len(ctx.enabledDeviceExtensions)),
		PpEnabledExtensionNames: ctx.enabledDeviceExtensions,
	}

	var device vulkan.Device
	result := vulkan.CreateDevice(ctx.gpu.ref, &deviceCreateInfo, nil, &device)
	panicOnError(result, "create device instance")
	ctx.device = device

	var graphicsQueue vulkan.Queue
	vulkan.GetDeviceQueue(ctx.device, ctx.gpu.queueFamilies.graphicsIndex, 0, &graphicsQueue)
	ctx.graphicsQueue = graphicsQueue

	var presentQueue vulkan.Queue
	vulkan.GetDeviceQueue(ctx.device, ctx.gpu.queueFamilies.presentIndex, 0, &presentQueue)
	ctx.presentQueue = presentQueue
}

func (ctx *Context) createSwapchain() {
	log.DebugCore("Creating Vulkan swapchain")

	ctx.swapchainImageCount = determineImageCount(
		ctx.swapchainImageCount, ctx.surfaceCapabilities.MinImageCount, ctx.surfaceCapabilities.MaxImageCount,
	)

	var swapchainImageExtent vulkan.Extent2D
	if ctx.surfaceCapabilities.CurrentExtent.Width != vulkan.MaxUint32 {
		swapchainImageExtent.Width = ctx.surfaceCapabilities.CurrentExtent.Width
		swapchainImageExtent.Height = ctx.surfaceCapabilities.CurrentExtent.Height
	} else {
		width, height := ctx.nativeWindow.GetSize()
		swapchainImageExtent.Width = uint32(width)
		swapchainImageExtent.Height = uint32(height)
	}

	// Attempt to use Mailbox present mode if available, otherwise use FIFO
	// THIS BEHAVIOUR ENABLES VSYNC BY DEFAULT! Use PresentModeImmediate to
	// support disabled VSYNC.
	presentMode := vulkan.PresentModeFifo

	var presentModeCount uint32
	result := vulkan.GetPhysicalDeviceSurfacePresentModes(ctx.gpu.ref, ctx.surface, &presentModeCount, nil)
	panicOnError(result, "retrieve supported present modes")
	supportedPresentModes := make([]vulkan.PresentMode, presentModeCount)
	result = vulkan.GetPhysicalDeviceSurfacePresentModes(ctx.gpu.ref, ctx.surface, &presentModeCount, supportedPresentModes)
	panicOnError(result, "retrieve supported present modes")

	for _, supportedMode := range supportedPresentModes {
		if supportedMode == vulkan.PresentModeMailbox {
			presentMode = supportedMode
		}
	}

	swapchainCreateInfo := vulkan.SwapchainCreateInfo{
		SType:                 vulkan.StructureTypeSwapchainCreateInfo,
		Surface:               ctx.surface,
		MinImageCount:         ctx.swapchainImageCount,
		ImageFormat:           ctx.surfaceFormat.Format,
		ImageColorSpace:       ctx.surfaceFormat.ColorSpace,
		ImageExtent:           swapchainImageExtent,
		ImageArrayLayers:      1, // No stereoscopic rendering, which requires 2
		ImageUsage:            vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		ImageSharingMode:      vulkan.SharingModeExclusive,
		QueueFamilyIndexCount: 0,   // Ignored since sharing mode is exclusive
		PQueueFamilyIndices:   nil, // Ignored since sharing mode is exclusive
		PreTransform:          vulkan.SurfaceTransformIdentityBit,
		CompositeAlpha:        vulkan.CompositeAlphaOpaqueBit,
		PresentMode:           presentMode,
		Clipped:               vulkan.True,
		OldSwapchain:          nil,
	}
	var swapchain vulkan.Swapchain
	result = vulkan.CreateSwapchain(ctx.device, &swapchainCreateInfo, nil, &swapchain)
	panicOnError(result, "create swapchain")

	ctx.swapchain = swapchain

	var swapchainImageCount uint32
	result = vulkan.GetSwapchainImages(ctx.device, ctx.swapchain, &swapchainImageCount, nil)
	panicOnError(result, "retrieve swapchain image count")

	ctx.swapchainImageCount = swapchainImageCount
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
			Format:     ctx.surfaceFormat.Format,
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

func (ctx *Context) createDepthStencilImage() {
	log.DebugCore("Creating Vulkan depth stencil image")

	// Take the first supported format of the following formats
	desiredFormats := []vulkan.Format{
		vulkan.FormatD32SfloatS8Uint,
		vulkan.FormatD24UnormS8Uint,
		vulkan.FormatD16UnormS8Uint,
		vulkan.FormatD32Sfloat,
		vulkan.FormatD16Unorm,
	}
	for _, format := range desiredFormats {
		var formatProps vulkan.FormatProperties
		vulkan.GetPhysicalDeviceFormatProperties(ctx.gpu.ref, format, &formatProps)
		formatProps.Deref()

		if formatProps.OptimalTilingFeatures&vulkan.FormatFeatureFlags(vulkan.FormatFeatureDepthStencilAttachmentBit) != 0 {
			ctx.depthStencilFormat = format
			break
		}
	}

	if ctx.depthStencilFormat == vulkan.FormatUndefined {
		log.PanicCore("None of the desired depth stencil formats are supported")
	}

	// Check whether stencil is available
	ctx.stencilAvailable = ctx.depthStencilFormat == vulkan.FormatD32SfloatS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD24UnormS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD16UnormS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD32Sfloat

	imageCreateInfo := vulkan.ImageCreateInfo{
		SType:     vulkan.StructureTypeImageCreateInfo,
		Flags:     0,
		ImageType: vulkan.ImageType2d,
		Format:    ctx.depthStencilFormat,
		Extent: vulkan.Extent3D{
			Width:  ctx.surfaceCapabilities.CurrentExtent.Width,
			Height: ctx.surfaceCapabilities.CurrentExtent.Height,
			Depth:  1,
		},
		MipLevels:             1,
		ArrayLayers:           1,
		Samples:               vulkan.SampleCount1Bit,
		Tiling:                vulkan.ImageTilingOptimal,
		Usage:                 vulkan.ImageUsageFlags(vulkan.ImageUsageDepthStencilAttachmentBit),
		SharingMode:           vulkan.SharingModeExclusive,
		QueueFamilyIndexCount: 0,   // Ignored because of exclusive mode
		PQueueFamilyIndices:   nil, // Ignored because of exclusive mode
		InitialLayout:         vulkan.ImageLayoutUndefined,
	}

	var depthStencilImage vulkan.Image
	result := vulkan.CreateImage(ctx.device, &imageCreateInfo, nil, &depthStencilImage)
	panicOnError(result, "create depth stencil image")
	ctx.depthStencilImage = depthStencilImage

	var imageMemoryRequirements vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(ctx.device, ctx.depthStencilImage, &imageMemoryRequirements)
	imageMemoryRequirements.Deref()

	memoryTypeIndex := ctx.findMemoryTypeIndex(&imageMemoryRequirements, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit))
	if memoryTypeIndex == vulkan.MaxUint32 {
		log.PanicCore("Could not find memory type to allocate depth stencil image memory")
	}

	memoryAllocateInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  imageMemoryRequirements.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}
	var depthStencilImageMemory vulkan.DeviceMemory
	vulkan.AllocateMemory(ctx.device, &memoryAllocateInfo, nil, &depthStencilImageMemory)
	ctx.depthStencilImageMemory = depthStencilImageMemory
	vulkan.BindImageMemory(ctx.device, ctx.depthStencilImage, ctx.depthStencilImageMemory, 0)

	aspectMask := vulkan.ImageAspectDepthBit
	if ctx.stencilAvailable {
		aspectMask |= vulkan.ImageAspectStencilBit
	}

	imageViewCreateInfo := vulkan.ImageViewCreateInfo{
		SType:      vulkan.StructureTypeImageViewCreateInfo,
		Image:      ctx.depthStencilImage,
		ViewType:   vulkan.ImageViewType2d,
		Format:     ctx.depthStencilFormat,
		Components: vulkan.ComponentMapping{}, // Use identity mapping for rgba components
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectFlags(aspectMask),
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var depthStencilImageView vulkan.ImageView
	result = vulkan.CreateImageView(ctx.device, &imageViewCreateInfo, nil, &depthStencilImageView)
	panicOnError(result, "create depth stencil image view")
	ctx.depthStencilImageView = depthStencilImageView
}

func (ctx *Context) destroyDepthStencilImage() {
	vulkan.DestroyImageView(ctx.device, ctx.depthStencilImageView, nil)
	vulkan.FreeMemory(ctx.device, ctx.depthStencilImageMemory, nil)
	vulkan.DestroyImage(ctx.device, ctx.depthStencilImage, nil)
}

func (ctx *Context) createRenderPass() {
	attachments := make([]vulkan.AttachmentDescription, 2)

	// Depth attachment
	attachments[0] = vulkan.AttachmentDescription{
		Format:         ctx.depthStencilFormat,
		Samples:        vulkan.SampleCount1Bit,
		LoadOp:         vulkan.AttachmentLoadOpClear,
		StoreOp:        vulkan.AttachmentStoreOpDontCare,
		StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpStore,
		InitialLayout:  vulkan.ImageLayoutUndefined,
		FinalLayout:    vulkan.ImageLayoutDepthStencilAttachmentOptimal,
	}
	// Color attachment
	attachments[1] = vulkan.AttachmentDescription{
		Format:        ctx.surfaceFormat.Format,
		Samples:       vulkan.SampleCount1Bit,
		LoadOp:        vulkan.AttachmentLoadOpClear,
		StoreOp:       vulkan.AttachmentStoreOpStore,
		InitialLayout: vulkan.ImageLayoutUndefined,
		FinalLayout:   vulkan.ImageLayoutPresentSrc,
	}

	depthStencilAttachment := vulkan.AttachmentReference{
		Attachment: 0,
		Layout:     vulkan.ImageLayoutDepthStencilAttachmentOptimal,
	}
	colorAttachments := make([]vulkan.AttachmentReference, 1)
	colorAttachments[0] = vulkan.AttachmentReference{
		Attachment: 1, // Reference to the color attachment index
		Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
	}

	subPasses := make([]vulkan.SubpassDescription, 1)
	subPasses[0] = vulkan.SubpassDescription{
		PipelineBindPoint:       vulkan.PipelineBindPointGraphics,
		InputAttachmentCount:    0,   // No other sub passes to reference
		PInputAttachments:       nil, // No other sub passes to reference
		ColorAttachmentCount:    uint32(len(colorAttachments)),
		PColorAttachments:       colorAttachments,
		PDepthStencilAttachment: &depthStencilAttachment,
	}

	renderPassCreateInfo := vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(attachments)),
		PAttachments:    attachments,
		SubpassCount:    uint32(len(subPasses)),
		PSubpasses:      subPasses,
		DependencyCount: 0,   // No dependencies between sub passes
		PDependencies:   nil, // No dependencies between sub passes
	}

	var renderPass vulkan.RenderPass
	result := vulkan.CreateRenderPass(ctx.device, &renderPassCreateInfo, nil, &renderPass)
	panicOnError(result, "create render pass")
	ctx.renderPass = renderPass
}

func (ctx *Context) createFramebuffers() {
	ctx.framebuffers = make([]vulkan.Framebuffer, ctx.swapchainImageCount)

	for i := uint32(0); i < ctx.swapchainImageCount; i++ {
		attachments := make([]vulkan.ImageView, 2)
		attachments[0] = ctx.depthStencilImageView
		attachments[1] = ctx.swapchainImageViews[i]

		framebufferCreateInfo := vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      ctx.renderPass,
			AttachmentCount: uint32(len(attachments)),
			PAttachments:    attachments,
			Width:           ctx.surfaceCapabilities.CurrentExtent.Width,
			Height:          ctx.surfaceCapabilities.CurrentExtent.Height,
			Layers:          1,
		}
		var framebuffer vulkan.Framebuffer
		result := vulkan.CreateFramebuffer(ctx.device, &framebufferCreateInfo, nil, &framebuffer)
		panicOnError(result, "create framebuffer for swapchain image "+string(i))
		ctx.framebuffers[i] = framebuffer
	}
}

func (ctx *Context) destroyFramebuffers() {
	for _, framebuffer := range ctx.framebuffers {
		vulkan.DestroyFramebuffer(ctx.device, framebuffer, nil)
	}
}

func (ctx *Context) newFence() vulkan.Fence {
	fenceCreateInfo := vulkan.FenceCreateInfo{
		SType: vulkan.StructureTypeFenceCreateInfo,
	}

	var fence vulkan.Fence
	result := vulkan.CreateFence(ctx.device, &fenceCreateInfo, nil, &fence)
	panicOnError(result, "create fence")

	return fence
}

func (ctx *Context) createSynchronizations() {
	ctx.swapchainImageAvailable = ctx.newFence()
}

func (ctx *Context) getActiveFramebuffer() vulkan.Framebuffer {
	return ctx.framebuffers[ctx.activeSwapchainImageindex]
}

func (ctx *Context) getSurfaceSize() vulkan.Extent2D {
	return vulkan.Extent2D{
		Width:  ctx.surfaceCapabilities.CurrentExtent.Width,
		Height: ctx.surfaceCapabilities.CurrentExtent.Height,
	}
}

func (ctx *Context) beginRender() {
	var activeSwapchainImage uint32
	result := vulkan.AcquireNextImage(ctx.device, ctx.swapchain, vulkan.MaxUint64, nil, ctx.swapchainImageAvailable, &activeSwapchainImage)
	panicOnError(result, "retrieve active swapchain image index")
	ctx.activeSwapchainImageindex = activeSwapchainImage

	result = vulkan.WaitForFences(ctx.device, 1, []vulkan.Fence{ctx.swapchainImageAvailable}, vulkan.True, vulkan.MaxUint64)
	panicOnError(result, "wait for swapchain image fence")

	result = vulkan.ResetFences(ctx.device, 1, []vulkan.Fence{ctx.swapchainImageAvailable})
	panicOnError(result, "reset swapchain image fence")

	//result = vulkan.QueueWaitIdle(ctx.graphicsQueue)
	//panicOnError(result, "wait for graphics queue to be idle")
}

func (ctx *Context) endRender(waitSemaphores []vulkan.Semaphore) {
	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: uint32(len(waitSemaphores)),
		PWaitSemaphores:    waitSemaphores,
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{ctx.swapchain},
		PImageIndices:      []uint32{ctx.activeSwapchainImageindex},
	}
	result := vulkan.QueuePresent(ctx.graphicsQueue, &presentInfo)
	panicOnError(result, "issue queue operations (draw call)")
}

func (ctx *Context) createCommandPool() {
	commandPoolCreateInfo := vulkan.CommandPoolCreateInfo{
		SType:            vulkan.StructureTypeCommandPoolCreateInfo,
		Flags:            vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateTransientBit | vulkan.CommandPoolCreateResetCommandBufferBit),
		QueueFamilyIndex: ctx.gpu.queueFamilies.graphicsIndex,
	}

	var commandPool vulkan.CommandPool
	result := vulkan.CreateCommandPool(ctx.device, &commandPoolCreateInfo, nil, &commandPool)
	panicOnError(result, "create command pool")
	ctx.commandPool = commandPool
}

func (ctx *Context) createCommandBuffer() {
	commandBufferAllocateInfo := vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        ctx.commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	commandBuffers := make([]vulkan.CommandBuffer, 1)
	result := vulkan.AllocateCommandBuffers(ctx.device, &commandBufferAllocateInfo, commandBuffers)
	panicOnError(result, "allocate command buffer")
	ctx.commandBuffer = commandBuffers[0]
}

func (ctx *Context) createRenderCompleteSemaphore() {
	semaphoreCreateInfo := vulkan.SemaphoreCreateInfo{
		SType: vulkan.StructureTypeSemaphoreCreateInfo,
	}

	var renderCompleteSemaphore vulkan.Semaphore
	result := vulkan.CreateSemaphore(ctx.device, &semaphoreCreateInfo, nil, &renderCompleteSemaphore)
	panicOnError(result, "create render complete semaphore")
	ctx.renderCompleteSemaphore = renderCompleteSemaphore
}

func (ctx *Context) Render() {
	ctx.beginRender()

	commandBufferBeginInfo := vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: vulkan.CommandBufferUsageFlags(vulkan.CommandBufferUsageOneTimeSubmitBit),
	}
	vulkan.BeginCommandBuffer(ctx.commandBuffer, &commandBufferBeginInfo)

	renderArea := vulkan.Rect2D{
		Offset: vulkan.Offset2D{
			X: 0,
			Y: 0,
		},
		Extent: ctx.getSurfaceSize(),
	}

	clearValues := make([]vulkan.ClearValue, 2)
	clearValues[0].SetDepthStencil(0.0, 0)
	clearValues[1].SetColor([]float32{0.8, 0.2, 0.2, 1.0})

	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:           vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:      ctx.renderPass,
		Framebuffer:     ctx.getActiveFramebuffer(),
		RenderArea:      renderArea,
		ClearValueCount: uint32(len(clearValues)),
		PClearValues:    clearValues,
	}
	vulkan.CmdBeginRenderPass(ctx.commandBuffer, &renderPassBeginInfo, vulkan.SubpassContentsInline)
	vulkan.CmdEndRenderPass(ctx.commandBuffer)

	vulkan.EndCommandBuffer(ctx.commandBuffer)

	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   0,
		PWaitSemaphores:      nil,
		PWaitDstStageMask:    nil,
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{ctx.commandBuffer},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{ctx.renderCompleteSemaphore},
	}
	vulkan.QueueSubmit(ctx.graphicsQueue, 1, []vulkan.SubmitInfo{submitInfo}, nil)

	ctx.endRender([]vulkan.Semaphore{ctx.renderCompleteSemaphore})
}
