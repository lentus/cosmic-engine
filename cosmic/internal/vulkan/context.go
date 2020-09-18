package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
	"strconv"
	"time"
)

const maxFramesInFlight = 2

type Context struct {
	nativeWindow *glfw.Window

	instance vulkan.Instance

	surface       surface
	gpu           physicalDevice
	device        vulkan.Device
	graphicsQueue vulkan.Queue
	presentQueue  vulkan.Queue

	swapchain           vulkan.Swapchain
	swapchainImageCount uint32
	swapchainImages     []vulkan.Image
	swapchainImageViews []vulkan.ImageView

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

	commandPool    vulkan.CommandPool
	commandBuffers []vulkan.CommandBuffer

	imageAvailableSemaphores []vulkan.Semaphore
	renderCompleteSemaphores []vulkan.Semaphore
	frameInFlightFences      []vulkan.Fence
	imagesInFlightFences     []vulkan.Fence
	currentFrame             int
}

func NewContext(nativeWindow *glfw.Window) *Context {
	log.InfoCore("Creating Vulkan graphics context")

	if !glfw.VulkanSupported() {
		log.PanicCore("glfw reports that Vulkan is not supported, aborting")
	}

	ctx := Context{
		nativeWindow:              nativeWindow,
		enabledInstanceLayers:     make([]string, 0),
		enabledInstanceExtensions: make([]string, 0),
		enabledDeviceExtensions:   make([]string, 0),
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
	ctx.createSwapchain()
	ctx.createSwapchainImages()

	// Graphics pipeline
	//ctx.createDepthStencilImage()
	ctx.createRenderPass()
	ctx.createFramebuffers()

	ctx.createCommandPool()
	ctx.createCommandBuffers()
	ctx.createSynchronizations()

	return &ctx
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.DeviceWaitIdle(ctx.device) // Wait for the graphics queue to be idle

	ctx.destroySynchronizations()
	vulkan.DestroyCommandPool(ctx.device, ctx.commandPool, nil)

	ctx.destroyFramebuffers()
	vulkan.DestroyRenderPass(ctx.device, ctx.renderPass, nil)
	//ctx.destroyDepthStencilImage()
	ctx.destroySwapchainImageViews()
	vulkan.DestroySwapchain(ctx.device, ctx.swapchain, nil) // Destroys swapchain images as well
	vulkan.DestroySurface(ctx.instance, ctx.surface.ref, nil)
	vulkan.DestroyDevice(ctx.device, nil)
	vulkan.DestroyDebugReportCallback(ctx.instance, ctx.debugCallback, nil)
	vulkan.DestroyInstance(ctx.instance, nil)
}

func (ctx *Context) createVulkanInstance() {
	log.DebugCore("Creating Vulkan instance")

	// TODO get version info and application name from somewhere
	applicationInfo := vulkan.ApplicationInfo{
		SType:         vulkan.StructureTypeApplicationInfo,
		ApiVersion:    vulkan.ApiVersion10,
		PEngineName:   safeStr("Cosmic Engine"),
		EngineVersion: vulkan.MakeVersion(0, 1, 0),
	}

	instanceCreateInfo := vulkan.InstanceCreateInfo{
		SType:                   vulkan.StructureTypeInstanceCreateInfo,
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

func (ctx *Context) createLogicalDevice() {
	log.DebugCore("Creating Vulkan device")

	queueFamilyIndices := []uint32{ctx.gpu.queueFamilies.graphicsIndex}
	if ctx.gpu.queueFamilies.hasSeparatePresentQueue() {
		queueFamilyIndices = append(queueFamilyIndices, ctx.gpu.queueFamilies.presentIndex)
	}

	var queueCreateInfos []vulkan.DeviceQueueCreateInfo
	for _, queueFamilyIndex := range queueFamilyIndices {
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
			Width:  ctx.surface.capabilities.CurrentExtent.Width,
			Height: ctx.surface.capabilities.CurrentExtent.Height,
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
	attachments := make([]vulkan.AttachmentDescription, 1)

	// Color attachment
	attachments[0] = vulkan.AttachmentDescription{
		Format:         ctx.surface.format.Format,
		Samples:        vulkan.SampleCount1Bit,
		LoadOp:         vulkan.AttachmentLoadOpClear,
		StoreOp:        vulkan.AttachmentStoreOpStore,
		StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
		InitialLayout:  vulkan.ImageLayoutUndefined,
		FinalLayout:    vulkan.ImageLayoutPresentSrc,
	}

	colorAttachmentRefs := make([]vulkan.AttachmentReference, 1)
	colorAttachmentRefs[0] = vulkan.AttachmentReference{
		Attachment: 0, // Reference to the color attachment index
		Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
	}

	subPasses := make([]vulkan.SubpassDescription, 1)
	subPasses[0] = vulkan.SubpassDescription{
		PipelineBindPoint:    vulkan.PipelineBindPointGraphics,
		ColorAttachmentCount: uint32(len(colorAttachmentRefs)),
		PColorAttachments:    colorAttachmentRefs,
	}

	// Make sure the subpass is not processed before it can write to the color attachment
	subpassDependencies := make([]vulkan.SubpassDependency, 1)
	subpassDependencies[0] = vulkan.SubpassDependency{
		SrcSubpass:    vulkan.SubpassExternal,
		DstSubpass:    0,
		SrcStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
		DstStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
		SrcAccessMask: 0,
		DstAccessMask: vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit),
	}

	renderPassCreateInfo := vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(attachments)),
		PAttachments:    attachments,
		SubpassCount:    uint32(len(subPasses)),
		PSubpasses:      subPasses,
		DependencyCount: 1,
		PDependencies:   subpassDependencies,
	}

	var renderPass vulkan.RenderPass
	result := vulkan.CreateRenderPass(ctx.device, &renderPassCreateInfo, nil, &renderPass)
	panicOnError(result, "create render pass")
	ctx.renderPass = renderPass
}

func (ctx *Context) createFramebuffers() {
	ctx.framebuffers = make([]vulkan.Framebuffer, ctx.swapchainImageCount)

	for i := range ctx.framebuffers {
		attachments := []vulkan.ImageView{ctx.swapchainImageViews[i]}

		framebufferCreateInfo := vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      ctx.renderPass,
			AttachmentCount: uint32(len(attachments)),
			PAttachments:    attachments,
			Width:           ctx.surface.capabilities.CurrentExtent.Width,
			Height:          ctx.surface.capabilities.CurrentExtent.Height,
			Layers:          1,
		}

		var framebuffer vulkan.Framebuffer
		result := vulkan.CreateFramebuffer(ctx.device, &framebufferCreateInfo, nil, &framebuffer)
		panicOnError(result, "create framebuffer for swapchain image "+strconv.Itoa(i))
		ctx.framebuffers[i] = framebuffer
	}
}

func (ctx *Context) destroyFramebuffers() {
	for _, framebuffer := range ctx.framebuffers {
		vulkan.DestroyFramebuffer(ctx.device, framebuffer, nil)
	}
}

func (ctx *Context) createCommandPool() {
	commandPoolCreateInfo := vulkan.CommandPoolCreateInfo{
		SType:            vulkan.StructureTypeCommandPoolCreateInfo,
		QueueFamilyIndex: ctx.gpu.queueFamilies.graphicsIndex,
		Flags:            0,
	}

	var commandPool vulkan.CommandPool
	result := vulkan.CreateCommandPool(ctx.device, &commandPoolCreateInfo, nil, &commandPool)
	panicOnError(result, "create command pool")
	ctx.commandPool = commandPool
}

func (ctx *Context) createCommandBuffers() {
	ctx.commandBuffers = make([]vulkan.CommandBuffer, ctx.swapchainImageCount)
	commandBufferAllocateInfo := vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        ctx.commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: uint32(len(ctx.commandBuffers)),
	}

	result := vulkan.AllocateCommandBuffers(ctx.device, &commandBufferAllocateInfo, ctx.commandBuffers)
	panicOnError(result, "allocate command buffers")

	// Record command buffers
	for i := range ctx.commandBuffers {
		beginInfo := vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
		}
		result = vulkan.BeginCommandBuffer(ctx.commandBuffers[i], &beginInfo)
		panicOnError(result, "start recording command buffer "+strconv.Itoa(i))

		renderArea := vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: ctx.getSurfaceSize(),
		}

		clearValues := make([]vulkan.ClearValue, 1)
		clearValues[0].SetColor([]float32{0.8, 0.2, 0.2, 1.0})

		renderPassBeginInfo := vulkan.RenderPassBeginInfo{
			SType:           vulkan.StructureTypeRenderPassBeginInfo,
			RenderPass:      ctx.renderPass,
			Framebuffer:     ctx.framebuffers[i],
			RenderArea:      renderArea,
			ClearValueCount: uint32(len(clearValues)),
			PClearValues:    clearValues,
		}

		vulkan.CmdBeginRenderPass(ctx.commandBuffers[i], &renderPassBeginInfo, vulkan.SubpassContentsInline)
		vulkan.CmdEndRenderPass(ctx.commandBuffers[i])

		result = vulkan.EndCommandBuffer(ctx.commandBuffers[i])
		panicOnError(result, "stop recording command buffer "+strconv.Itoa(i))
	}
}

func (ctx *Context) getSurfaceSize() vulkan.Extent2D {
	return vulkan.Extent2D{
		Width:  ctx.surface.capabilities.CurrentExtent.Width,
		Height: ctx.surface.capabilities.CurrentExtent.Height,
	}
}

func (ctx *Context) createSynchronizations() {
	ctx.imageAvailableSemaphores = make([]vulkan.Semaphore, maxFramesInFlight)
	ctx.renderCompleteSemaphores = make([]vulkan.Semaphore, maxFramesInFlight)
	ctx.frameInFlightFences = make([]vulkan.Fence, maxFramesInFlight)

	for i := range ctx.imageAvailableSemaphores {
		ctx.imageAvailableSemaphores[i] = ctx.newSemaphore()
		ctx.renderCompleteSemaphores[i] = ctx.newSemaphore()
		ctx.frameInFlightFences[i] = ctx.newFence()
	}

	ctx.imagesInFlightFences = make([]vulkan.Fence, ctx.swapchainImageCount)
}

func (ctx *Context) destroySynchronizations() {
	for i := range ctx.frameInFlightFences {
		vulkan.DestroyFence(ctx.device, ctx.frameInFlightFences[i], nil)
		vulkan.DestroySemaphore(ctx.device, ctx.renderCompleteSemaphores[i], nil)
		vulkan.DestroySemaphore(ctx.device, ctx.imageAvailableSemaphores[i], nil)
	}
}

func (ctx *Context) newFence() vulkan.Fence {
	fenceCreateInfo := vulkan.FenceCreateInfo{
		SType: vulkan.StructureTypeFenceCreateInfo,
		Flags: vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit),
	}

	var fence vulkan.Fence
	result := vulkan.CreateFence(ctx.device, &fenceCreateInfo, nil, &fence)
	panicOnError(result, "create fence")

	return fence
}

func (ctx *Context) newSemaphore() vulkan.Semaphore {
	semaphoreCreateInfo := vulkan.SemaphoreCreateInfo{
		SType: vulkan.StructureTypeSemaphoreCreateInfo,
	}

	var semaphore vulkan.Semaphore
	result := vulkan.CreateSemaphore(ctx.device, &semaphoreCreateInfo, nil, &semaphore)
	panicOnError(result, "create semaphore")
	return semaphore
}

func (ctx *Context) Render() {
	timeout := uint64(5 * time.Millisecond.Nanoseconds())

	// Wait for frame to be presented if still in flight
	result := vulkan.WaitForFences(
		ctx.device, 1, []vulkan.Fence{ctx.frameInFlightFences[ctx.currentFrame]}, vulkan.True, timeout,
	)
	if result != vulkan.Success {
		log.PanicfCore("%s while waiting for frame in flight fence %d", fmtResult(result), ctx.currentFrame)
	}

	var imageIndex uint32
	result = vulkan.AcquireNextImage(
		ctx.device, ctx.swapchain, vulkan.MaxUint64,
		ctx.imageAvailableSemaphores[ctx.currentFrame],
		vulkan.NullFence, &imageIndex,
	)
	panicOnError(result, "retrieve active swapchain image index")

	// Check whether the swapchain image is currently in flight. After the first
	// ctx.swpapchainImageCount frames these will always be filled, but waiting
	// for fences that were already signalled is just a no-op.
	if ctx.imagesInFlightFences[imageIndex] != nil {
		result = vulkan.WaitForFences(
			ctx.device, 1, []vulkan.Fence{ctx.imagesInFlightFences[imageIndex]}, vulkan.True, timeout,
		)
		if result != vulkan.Success {
			log.PanicfCore("%s while waiting for image in flight fence %d", fmtResult(result), ctx.currentFrame)
		}
	}
	ctx.imagesInFlightFences[imageIndex] = ctx.frameInFlightFences[ctx.currentFrame]

	pipelineStageFlags := vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)
	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{ctx.imageAvailableSemaphores[ctx.currentFrame]},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{pipelineStageFlags},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{ctx.commandBuffers[imageIndex]},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{ctx.renderCompleteSemaphores[ctx.currentFrame]},
	}
	vulkan.ResetFences(ctx.device, 1, []vulkan.Fence{ctx.frameInFlightFences[ctx.currentFrame]})
	result = vulkan.QueueSubmit(
		ctx.graphicsQueue, 1, []vulkan.SubmitInfo{submitInfo}, ctx.frameInFlightFences[ctx.currentFrame],
	)
	panicOnError(result, "submit draw command buffer")

	// TODO DEBUG wait for command buffer to finish execution
	result = vulkan.WaitForFences(
		ctx.device, 1, []vulkan.Fence{ctx.frameInFlightFences[ctx.currentFrame]}, vulkan.True, timeout,
	)
	if result != vulkan.Success {
		log.PanicfCore("%s while waiting for frame in flight fence %d", fmtResult(result), ctx.currentFrame)
	}

	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{ctx.renderCompleteSemaphores[ctx.currentFrame]},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{ctx.swapchain},
		PImageIndices:      []uint32{imageIndex},
	}
	result = vulkan.QueuePresent(ctx.presentQueue, &presentInfo)
	panicOnError(result, "queue present")

	log.DebugfCore("Finished draw call (swapchain image %d, frame %d)", imageIndex, ctx.currentFrame)
	ctx.currentFrame = (ctx.currentFrame + 1) % maxFramesInFlight
}
