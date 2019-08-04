package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type Context struct {
	nativeWindow        *glfw.Window
	instance            vulkan.Instance
	gpu                 vulkan.PhysicalDevice
	device              vulkan.Device
	graphicsFamilyIndex uint32
}

func NewContext(nativeWindow *glfw.Window) *Context {
	log.InfoCore("Creating Vulkan graphics context")
	ctx := Context{nativeWindow: nativeWindow}
	ctx.nativeWindow.MakeContextCurrent()

	vulkan.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vulkan.Init(); err != nil {
		log.PanicfCore("failed to initialise Vulkan, %s", err.Error())
	}

	ctx.createVulkanInstance()
	ctx.selectGpu()
	ctx.findGraphicsQueueFamily()

	// Setup debugging
	ctx.getInstanceLayers()
	ctx.getDeviceLayers()

	ctx.createDevice()

	return &ctx
}

func (ctx *Context) SwapBuffers() {
	ctx.nativeWindow.SwapBuffers()
}

func (ctx *Context) Terminate() {
	log.DebugCore("Terminating Vulkan graphics context")
	vulkan.DestroyDevice(ctx.device, nil)
	vulkan.DestroyInstance(ctx.instance, nil)
}

func (ctx *Context) createVulkanInstance() {
	log.DebugCore("Creating Vulkan instance")

	// TODO get version info and application name from somewhere
	applicationInfo := vulkan.ApplicationInfo{
		SType:         vulkan.StructureTypeApplicationInfo,
		ApiVersion:    vulkan.MakeVersion(1, 1, 88),
		PEngineName:   "Cosmic Engine\x00",
		EngineVersion: vulkan.MakeVersion(0, 1, 0),
	}

	instanceCreateInfo := vulkan.InstanceCreateInfo{
		SType:            vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &applicationInfo,
	}

	var instance vulkan.Instance
	result := vulkan.CreateInstance(&instanceCreateInfo, nil, &instance)
	panicOnError(result, "create Vulkan instance")

	ctx.instance = instance
	//vulkan.InitInstance(ctx.instance)
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

func (ctx *Context) getInstanceLayers() {
	var layerCount uint32
	vulkan.EnumerateInstanceLayerProperties(&layerCount, nil)
	layerPropertiesList := make([]vulkan.LayerProperties, layerCount)
	result := vulkan.EnumerateInstanceLayerProperties(&layerCount, layerPropertiesList)
	panicOnError(result, "retrieve instance layers")

	log.DebugfCore("Instance layers (%d):", len(layerPropertiesList))
	for _, props := range layerPropertiesList {
		props.Deref()
		log.DebugfCore("\t%s [%s]", props.LayerName, props.Description)
	}
}

func (ctx *Context) getDeviceLayers() {
	var layerCount uint32
	vulkan.EnumerateDeviceLayerProperties(ctx.gpu, &layerCount, nil)
	layerPropertiesList := make([]vulkan.LayerProperties, layerCount)
	result := vulkan.EnumerateDeviceLayerProperties(ctx.gpu, &layerCount, layerPropertiesList)
	panicOnError(result, "retrieve device layers")

	log.DebugfCore("Device layers (%d):", len(layerPropertiesList))
	for _, props := range layerPropertiesList {
		props.Deref()
		log.DebugfCore("\t%s [%s]", props.LayerName, props.Description)
	}
}
