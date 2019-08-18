package vulkan

import (
	"bytes"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
	"unsafe"
)

func (ctx *Context) setupDebug() {
	ctx.availableInstanceLayers = ctx.getInstanceLayers()

	// Setup debug layers
	ctx.enableIfAvailable("VK_LAYER_LUNARG_core_validation")
	ctx.enableIfAvailable("VK_LAYER_LUNARG_parameter_validation")
	ctx.enableIfAvailable("VK_LAYER_LUNARG_object_tracker")
	ctx.enableIfAvailable("VK_LAYER_GOOGLE_threading")
	ctx.enableIfAvailable("VK_LAYER_LUNARG_screenshot")
	ctx.enableIfAvailable("VK_LAYER_LUNARG_monitor")

	// Setup debug error callback
	ctx.enabledInstanceExtensions = append(ctx.enabledInstanceExtensions, safeStr(vulkan.ExtDebugReportExtensionName))
}

func (ctx *Context) enableIfAvailable(layerName string) {
	for _, instanceLayer := range ctx.availableInstanceLayers {
		instanceLayer.Deref()

		instanceLayerName := string(bytes.Trim(instanceLayer.LayerName[:], "\x00"))
		if instanceLayerName == layerName {
			log.DebugfCore("Enabling instance layer %s", instanceLayer.LayerName)
			ctx.enabledInstanceLayers = append(ctx.enabledInstanceLayers, string(instanceLayer.LayerName[:]))
			return
		}
	}

	log.WarnfCore("Cannot enable instance layer %s (not available)", layerName)
}

func (ctx *Context) getInstanceLayers() []vulkan.LayerProperties {
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

	return layerPropertiesList
}

func (ctx *Context) getDeviceLayers() []vulkan.LayerProperties {
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

	log.WarnCore("Device layers are deprecated since vulkan 1.0.13, and not supported by Cosmic")

	return layerPropertiesList
}

var debugBit = vulkan.DebugReportFlags(vulkan.DebugReportDebugBit)
var infoBit = vulkan.DebugReportFlags(vulkan.DebugReportInformationBit)
var warnBit = vulkan.DebugReportFlags(vulkan.DebugReportWarningBit)
var perfWarnBit = vulkan.DebugReportFlags(vulkan.DebugReportPerformanceWarningBit)
var errorBit = vulkan.DebugReportFlags(vulkan.DebugReportErrorBit)

func vulkanDebugReportCallback(
	flags vulkan.DebugReportFlags,
	objectType vulkan.DebugReportObjectType,
	_ uint64,
	_ uint,
	msgCode int32,
	layer string,
	msg string,
	_ unsafe.Pointer) vulkan.Bool32 {

	fmtString := "Vulkan %s [%s] - %s (code %d)"

	switch {
	case flags&infoBit != 0:
		log.InfofCore(fmtString, layer, fmtObjectType(objectType), msg, msgCode)
	case flags&warnBit != 0:
		log.WarnfCore(fmtString, layer, fmtObjectType(objectType), msg, msgCode)
	case flags&perfWarnBit != 0:
		log.WarnfCore(fmtString, "<Performance> "+layer, fmtObjectType(objectType), msg, msgCode)
	case flags&errorBit != 0:
		log.ErrorfCore(fmtString, layer, fmtObjectType(objectType), msg, msgCode)
	}

	return vulkan.False
}

func (ctx *Context) initDebug() {
	reportFlagBits := vulkan.DebugReportInformationBit |
		vulkan.DebugReportWarningBit |
		vulkan.DebugReportErrorBit |
		vulkan.DebugReportPerformanceWarningBit

	createInfo := vulkan.DebugReportCallbackCreateInfo{
		SType:       vulkan.StructureTypeDebugReportCallbackCreateInfo,
		Flags:       vulkan.DebugReportFlags(reportFlagBits),
		PfnCallback: vulkanDebugReportCallback,
	}

	var debugCallback vulkan.DebugReportCallback
	vulkan.CreateDebugReportCallback(ctx.instance, &createInfo, nil, &debugCallback)
	ctx.debugCallback = debugCallback
}

func (ctx *Context) deInitDebug() {
	vulkan.DestroyDebugReportCallback(ctx.instance, ctx.debugCallback, nil)
}
