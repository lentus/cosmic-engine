package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type Context struct {
	nativeWindow        *glfw.Window
	surface             vulkan.Surface
	surfaceFormat       vulkan.SurfaceFormat
	surfaceCapabilities vulkan.SurfaceCapabilities

	instance            vulkan.Instance
	gpu                 vulkan.PhysicalDevice
	device              vulkan.Device
	graphicsFamilyIndex uint32

	availableInstanceLayers     []vulkan.LayerProperties
	enabledInstanceLayers       []string
	availableInstanceExtensions []vulkan.ExtensionProperties
	enabledInstanceExtensions   []string

	debugCallback vulkan.DebugReportCallback
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
	}
	ctx.nativeWindow.MakeContextCurrent()

	vulkan.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vulkan.Init(); err != nil {
		log.PanicfCore("failed to initialise Vulkan, %s", err.Error())
	}

	ctx.availableInstanceLayers = ctx.getInstanceLayers()
	ctx.availableInstanceExtensions = ctx.getInstanceExtensions()

	ctx.setupDebug()
	ctx.createVulkanInstance()
	ctx.initDebug()
	ctx.selectGpu()
	ctx.findGraphicsQueueFamily()
	ctx.createDevice()
	ctx.createWindowSurface()

	return &ctx
}

func (ctx *Context) SwapBuffers() {
	ctx.nativeWindow.SwapBuffers()
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.DestroyDevice(ctx.device, nil)
	ctx.deInitDebug()
	vulkan.DestroyInstance(ctx.instance, nil)
}

func (ctx *Context) createVulkanInstance() {
	log.DebugCore("Creating Vulkan instance")

	ctx.setupLayersAndExtensions()

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
		SType:                vulkan.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount: 1,
		PQueueCreateInfos:    queueCreateInfos,
	}

	var device vulkan.Device
	result := vulkan.CreateDevice(ctx.gpu, &deviceCreateInfo, nil, &device)
	panicOnError(result, "create device instance")

	ctx.device = device
}

func (ctx *Context) createWindowSurface() {
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
