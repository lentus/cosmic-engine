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

	swapchain            vulkan.Swapchain
	swapchainImageCount  uint32
	swapchainImages      []vulkan.Image
	swapchainImageViews  []vulkan.ImageView
	swapchainImageExtent vulkan.Extent2D

	pipelineLayout   vulkan.PipelineLayout
	renderPass       vulkan.RenderPass
	graphicsPipeline vulkan.Pipeline

	framebuffers []vulkan.Framebuffer

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
	ctx.createRenderPass()
	ctx.createGraphicsPipeline()
	ctx.createFramebuffers()

	ctx.createCommandPool()
	ctx.createCommandBuffers()
	ctx.createSynchronizations()

	return &ctx
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.DeviceWaitIdle(ctx.device) // Wait for the device to be idle

	ctx.destroySynchronizations()
	vulkan.DestroyCommandPool(ctx.device, ctx.commandPool, nil)
	ctx.destroyFramebuffers()

	ctx.destroyGraphicsPipeline()

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
			Extent: vulkan.Extent2D{
				Width:  ctx.swapchainImageExtent.Width,
				Height: ctx.swapchainImageExtent.Height,
			},
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
		vulkan.CmdBindPipeline(ctx.commandBuffers[i], vulkan.PipelineBindPointGraphics, ctx.graphicsPipeline)
		vulkan.CmdDraw(ctx.commandBuffers[i], 3, 1, 0, 0)
		vulkan.CmdEndRenderPass(ctx.commandBuffers[i])

		result = vulkan.EndCommandBuffer(ctx.commandBuffers[i])
		panicOnError(result, "stop recording command buffer "+strconv.Itoa(i))
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
