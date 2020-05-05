package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/graphics"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type Context struct {
	nativeWindow *glfw.Window

	instance            vulkan.Instance
	gpu                 vulkan.PhysicalDevice
	device              vulkan.Device
	graphicsFamilyIndex uint32

	surface             vulkan.Surface
	surfaceFormat       vulkan.SurfaceFormat
	surfaceCapabilities vulkan.SurfaceCapabilities
	swapchain           vulkan.Swapchain
	swapchainImageCount uint32

	availableInstanceLayers     []vulkan.LayerProperties
	availableInstanceExtensions []vulkan.ExtensionProperties
	availableDeviceExtensions   []vulkan.ExtensionProperties
	enabledInstanceLayers       []string
	enabledInstanceExtensions   []string
	enabledDeviceExtensions     []string

	debugCallback vulkan.DebugReportCallback
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

	ctx.setupDebug()
	ctx.createVulkanInstance()
	ctx.initDebugCallback()
	ctx.selectGpu()
	ctx.findGraphicsQueueFamily()
	ctx.createDevice()
	ctx.createWindowSurface()
	ctx.createSwapchain()

	return &ctx
}

func (ctx *Context) SwapBuffers() {
	ctx.nativeWindow.SwapBuffers()
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.DestroySwapchain(ctx.device, ctx.swapchain, nil)
	vulkan.DestroySurface(ctx.instance, ctx.surface, nil)
	vulkan.DestroyDevice(ctx.device, nil)
	vulkan.DestroyDebugReportCallback(ctx.instance, ctx.debugCallback, nil)
	vulkan.DestroyInstance(ctx.instance, nil)
}

func (ctx *Context) createVulkanInstance() {
	log.DebugCore("Creating Vulkan instance")

	ctx.setupInstanceLayersAndExtensions()

	// TODO get version info and application name from somewhere
	applicationInfo := vulkan.ApplicationInfo{
		SType:         vulkan.StructureTypeApplicationInfo,
		ApiVersion:    vulkan.MakeVersion(1, 1, 88),
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

func (ctx *Context) selectGpu() {
	log.DebugCore("Selecting gpu")

	var gpuCount uint32
	vulkan.EnumeratePhysicalDevices(ctx.instance, &gpuCount, nil)
	gpus := make([]vulkan.PhysicalDevice, gpuCount)
	result := vulkan.EnumeratePhysicalDevices(ctx.instance, &gpuCount, gpus)
	panicOnError(result, "retrieve gpu list")

	log.DebugfCore("Found %d gpu(s)", len(gpus))
	for _, gpu := range gpus {
		var gpuProperties vulkan.PhysicalDeviceProperties
		vulkan.GetPhysicalDeviceProperties(gpu, &gpuProperties)
		gpuProperties.Deref()

		var memoryProperties vulkan.PhysicalDeviceMemoryProperties
		vulkan.GetPhysicalDeviceMemoryProperties(gpu, &memoryProperties)
		memoryProperties.Deref()

		log.DebugfCore("\tName           %s", string(gpuProperties.DeviceName[:]))
		log.DebugfCore("\tID             %d", gpuProperties.DeviceID)
		log.DebugfCore("\tType           %d", gpuProperties.DeviceType)
		log.DebugfCore("\tAPI version    %d", gpuProperties.ApiVersion)
		log.DebugfCore("\tVendor ID      %d", gpuProperties.VendorID)
		log.DebugfCore("\tDriver version %d", gpuProperties.DriverVersion)
	}

	// Select a GPU TODO find best performing one
	ctx.gpu = gpus[0]
}

func (ctx *Context) findGraphicsQueueFamily() {
	var familyCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(ctx.gpu, &familyCount, nil)
	queueFamiliePropertiesList := make([]vulkan.QueueFamilyProperties, familyCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(ctx.gpu, &familyCount, queueFamiliePropertiesList)

	found := false
	for i, properties := range queueFamiliePropertiesList {
		properties.Deref()
		if properties.QueueFlags&vulkan.QueueFlags(vulkan.QueueGraphicsBit) != 0 {
			ctx.graphicsFamilyIndex = uint32(i)
			found = true
			break
		}
	}

	if !found {
		log.PanicCore("Failed to find queue family supporting graphics on selected gpu")
	}
}

func (ctx *Context) createDevice() {
	log.DebugCore("Creating Vulkan device")

	ctx.setupDeviceExtensions()

	deviceQueueCreateInfo := vulkan.DeviceQueueCreateInfo{
		SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: ctx.graphicsFamilyIndex,
		QueueCount:       1,
		PQueuePriorities: []float32{1.0},
	}

	queueCreateInfos := []vulkan.DeviceQueueCreateInfo{
		deviceQueueCreateInfo,
	}

	deviceCreateInfo := vulkan.DeviceCreateInfo{
		SType:                   vulkan.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    1,
		PQueueCreateInfos:       queueCreateInfos,
		EnabledExtensionCount:   uint32(len(ctx.enabledDeviceExtensions)),
		PpEnabledExtensionNames: ctx.enabledDeviceExtensions,
	}

	var device vulkan.Device
	result := vulkan.CreateDevice(ctx.gpu, &deviceCreateInfo, nil, &device)
	panicOnError(result, "create device instance")

	ctx.device = device
}

func (ctx *Context) createWindowSurface() {
	log.DebugCore("Creating window surface")

	surfacePtr, err := ctx.nativeWindow.CreateWindowSurface(ctx.instance, nil)
	if err != nil {
		log.PanicfCore("failed to create vulkan window surface, %s", err.Error())
	}
	ctx.surface = vulkan.SurfaceFromPointer(surfacePtr)

	var wsiSupported vulkan.Bool32
	result := vulkan.GetPhysicalDeviceSurfaceSupport(ctx.gpu, ctx.graphicsFamilyIndex, ctx.surface, &wsiSupported)
	panicOnError(result, "check whether WSI is supported")

	if wsiSupported == vulkan.False {
		log.PanicCore("the GLFW surface does not support WSI")
	}

	surfaceCapabilities := vulkan.SurfaceCapabilities{}
	result = vulkan.GetPhysicalDeviceSurfaceCapabilities(ctx.gpu, ctx.surface, &surfaceCapabilities)
	panicOnError(result, "get surface capabilities")

	surfaceCapabilities.Deref()
	surfaceCapabilities.CurrentExtent.Deref()
	surfaceCapabilities.MinImageExtent.Deref()
	surfaceCapabilities.MaxImageExtent.Deref()
	ctx.surfaceCapabilities = surfaceCapabilities

	var formatCount uint32
	result = vulkan.GetPhysicalDeviceSurfaceFormats(ctx.gpu, ctx.surface, &formatCount, nil)
	panicOnError(result, "get physical device format count")
	if formatCount == 0 {
		log.PanicCore("no surface format found")
	}

	surfaceFormats := make([]vulkan.SurfaceFormat, formatCount)
	vulkan.GetPhysicalDeviceSurfaceFormats(ctx.gpu, ctx.surface, &formatCount, surfaceFormats)
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
	result := vulkan.GetPhysicalDeviceSurfacePresentModes(ctx.gpu, ctx.surface, &presentModeCount, nil)
	panicOnError(result, "retrieve supported present modes")
	supportedPresentModes := make([]vulkan.PresentMode, presentModeCount)
	result = vulkan.GetPhysicalDeviceSurfacePresentModes(ctx.gpu, ctx.surface, &presentModeCount, supportedPresentModes)
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
}

func determineImageCount(requested, min, max uint32) uint32 {
	if requested < min {
		log.WarnfCore("Requested image count %d not supported by your system, min is %d", requested, min)
		return min
	}

	if max > 0 && requested > max {
		log.WarnfCore("Requested image count %d not supported by your system, max is %d", requested, max)
		return max
	}

	return requested
}
