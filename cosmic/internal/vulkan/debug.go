package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
	"unsafe"
)

func (ctx *Context) setupDebug() {
	// Setup debug layers
	ctx.enableLayerIfAvailable("VK_LAYER_KHRONOS_validation")
	ctx.enableLayerIfAvailable("VK_LAYER_LUNARG_screenshot")
	ctx.enableLayerIfAvailable("VK_LAYER_LUNARG_monitor")

	// Setup debug error callback
	ctx.enabledInstanceExtensions = append(ctx.enabledInstanceExtensions, safeStr(vulkan.ExtDebugReportExtensionName))
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
	case flags&debugBit != 0:
		log.DebugfCore(fmtString, layer, fmtObjectType(objectType), msg, msgCode)
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

func (ctx *Context) initDebugCallback() {
	createInfo := createDebugReportCallbackCreateInfo()

	var debugCallback vulkan.DebugReportCallback
	vulkan.CreateDebugReportCallback(ctx.instance, &createInfo, nil, &debugCallback)
	ctx.debugCallback = debugCallback
}

func createDebugReportCallbackCreateInfo() vulkan.DebugReportCallbackCreateInfo {
	reportFlagBits :=
		//vulkan.DebugReportDebugBit |
		//	vulkan.DebugReportInformationBit |
		vulkan.DebugReportWarningBit |
			vulkan.DebugReportErrorBit |
			vulkan.DebugReportPerformanceWarningBit

	createInfo := vulkan.DebugReportCallbackCreateInfo{
		SType:       vulkan.StructureTypeDebugReportCallbackCreateInfo,
		Flags:       vulkan.DebugReportFlags(reportFlagBits),
		PfnCallback: vulkanDebugReportCallback,
	}

	return createInfo
}
