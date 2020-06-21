package vulkan

import (
	"bytes"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

type physicalDevice struct {
	ref              vulkan.PhysicalDevice
	properties       vulkan.PhysicalDeviceProperties
	memoryProperties vulkan.PhysicalDeviceMemoryProperties
	features         vulkan.PhysicalDeviceFeatures
	queueFamilies    queueFamilies
}

func (ctx *Context) selectPhysicalDevice() {
	log.DebugCore("Selecting gpu")

	var gpuCount uint32
	vulkan.EnumeratePhysicalDevices(ctx.instance, &gpuCount, nil)
	gpus := make([]vulkan.PhysicalDevice, gpuCount)
	result := vulkan.EnumeratePhysicalDevices(ctx.instance, &gpuCount, gpus)
	panicOnError(result, "retrieve gpu list")

	log.DebugfCore("Found %d gpu(s)", len(gpus))
	var selected vulkan.PhysicalDevice
	for _, gpu := range gpus {
		gpuProperties := getProperties(gpu)

		log.DebugfCore("\tName           %s", gpuProperties.DeviceName)
		log.DebugfCore("\tID             %d", gpuProperties.DeviceID)
		log.DebugfCore("\tType           %d", gpuProperties.DeviceType)
		log.DebugfCore("\tAPI version    %d", gpuProperties.ApiVersion)
		log.DebugfCore("\tVendor ID      %d", gpuProperties.VendorID)
		log.DebugfCore("\tDriver version %d", gpuProperties.DriverVersion)
		log.DebugCore("")

		if isDeviceSuitable(gpu, ctx.surface.ref, ctx.enabledDeviceExtensions) && selected == nil {
			// TODO find best performing one
			selected = gpu
		}
	}

	if selected == nil {
		log.PanicCore("failed to find a suitable gpu")
	}

	ctx.gpu = physicalDevice{
		ref:              selected,
		properties:       getProperties(selected),
		memoryProperties: getMemoryProperties(selected),
		features:         getFeatures(selected),
		queueFamilies:    findQueueFamilies(selected, ctx.surface.ref),
	}
	log.InfofCore("\tUsing %s", string(ctx.gpu.properties.DeviceName[:]))

	ctx.availableDeviceExtensions = getExtensions(ctx.gpu.ref)
	log.DebugfCore("Device extensions (%d):", len(ctx.availableDeviceExtensions))
	for _, deviceExtension := range ctx.availableDeviceExtensions {
		deviceExtension.Deref()
		log.DebugfCore("\t%s", deviceExtension.ExtensionName)
	}
}

func isDeviceSuitable(gpu vulkan.PhysicalDevice, surface vulkan.Surface, enabledExtensions []string) bool {
	// Check whether all required queue families are present
	indices := findQueueFamilies(gpu, surface)
	if !indices.complete() {
		return false
	}

	// Check whether all enabled device extensions are supported
	availableExtensions := getExtensions(gpu)
	availableExtensionMap := make(map[string]bool)
	for _, availableExtension := range availableExtensions {
		availableExtension.Deref()
		index := string(bytes.Trim(availableExtension.ExtensionName[:], "\x00"))
		availableExtensionMap[safeStr(index)] = true
	}

	for _, requiredExtension := range enabledExtensions {
		if _, extensionSupported := availableExtensionMap[requiredExtension]; !extensionSupported {
			return false
		}
	}

	// Check whether surface can be drawn to
	return len(getSurfaceFormats(gpu, surface)) > 0 && len(getPresentModes(gpu, surface)) > 0
}

func getProperties(gpu vulkan.PhysicalDevice) vulkan.PhysicalDeviceProperties {
	var gpuProperties vulkan.PhysicalDeviceProperties
	vulkan.GetPhysicalDeviceProperties(gpu, &gpuProperties)
	gpuProperties.Deref()

	return gpuProperties
}

func getMemoryProperties(gpu vulkan.PhysicalDevice) vulkan.PhysicalDeviceMemoryProperties {
	var memoryProperties vulkan.PhysicalDeviceMemoryProperties
	vulkan.GetPhysicalDeviceMemoryProperties(gpu, &memoryProperties)
	memoryProperties.Deref()

	return memoryProperties
}

func getFeatures(gpu vulkan.PhysicalDevice) vulkan.PhysicalDeviceFeatures {
	var gpuFeatures vulkan.PhysicalDeviceFeatures
	vulkan.GetPhysicalDeviceFeatures(gpu, &gpuFeatures)
	gpuFeatures.Deref()

	return gpuFeatures
}

func getExtensions(gpu vulkan.PhysicalDevice) []vulkan.ExtensionProperties {
	var extensionCount uint32
	vulkan.EnumerateDeviceExtensionProperties(gpu, "", &extensionCount, nil)
	extensionProperties := make([]vulkan.ExtensionProperties, extensionCount)
	vulkan.EnumerateDeviceExtensionProperties(gpu, "", &extensionCount, extensionProperties)

	return extensionProperties
}

type queueFamilies struct {
	graphicsIndex    uint32
	hasGraphicsIndex bool

	presentIndex    uint32
	hasPresentIndex bool
}

func (qf queueFamilies) complete() bool {
	return qf.hasGraphicsIndex && qf.hasPresentIndex
}

func (qf queueFamilies) hasSeparatePresentQueue() bool {
	return qf.graphicsIndex != qf.presentIndex
}

func findQueueFamilies(device vulkan.PhysicalDevice, surface vulkan.Surface) queueFamilies {
	var familyCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &familyCount, nil)
	queueFamiliePropertiesList := make([]vulkan.QueueFamilyProperties, familyCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device, &familyCount, queueFamiliePropertiesList)

	var queueFamilies queueFamilies
	for i, properties := range queueFamiliePropertiesList {
		properties.Deref()

		if properties.QueueFlags&vulkan.QueueFlags(vulkan.QueueGraphicsBit) != 0 {
			queueFamilies.graphicsIndex = uint32(i)
			queueFamilies.hasGraphicsIndex = true
		}

		var presentSupported vulkan.Bool32
		vulkan.GetPhysicalDeviceSurfaceSupport(device, uint32(i), surface, &presentSupported)
		if presentSupported == vulkan.True {
			queueFamilies.presentIndex = uint32(i)
			queueFamilies.hasPresentIndex = true
		}

		if queueFamilies.complete() {
			break
		}
	}

	return queueFamilies
}
