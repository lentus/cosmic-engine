package vulkan

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

func safeStr(str string) string {
	if len(str) == 0 {
		return "\x00"
	}

	if str[len(str)-1] != '\x00' {
		return str + "\x00"
	}

	return str
}

func panicOnError(result vulkan.Result, operation string) {
	if result < vulkan.Success {
		log.PanicfCore("failed to %s: %s", operation, fmtResult(result))
	}

	if result > vulkan.Success {
		log.WarnfCore("%s partially successful: %s", operation, fmtResult(result))
	}
}

func fmtResult(error vulkan.Result) string {
	var errName string

	switch error {
	case vulkan.NotReady:
		errName = "NotReady"
	case vulkan.Timeout:
		errName = "Timeout"
	case vulkan.EventSet:
		errName = "EventSet"
	case vulkan.EventReset:
		errName = "EventReset"
	case vulkan.Incomplete:
		errName = "Incomplete"
	case vulkan.ErrorOutOfHostMemory:
		errName = "ErrorOutOfHostMemory"
	case vulkan.ErrorOutOfDeviceMemory:
		errName = "ErrorOutOfDeviceMemory"
	case vulkan.ErrorInitializationFailed:
		errName = "ErrorInitializationFailed"
	case vulkan.ErrorDeviceLost:
		errName = "ErrorDeviceLost"
	case vulkan.ErrorMemoryMapFailed:
		errName = "ErrorMemoryMapFailed"
	case vulkan.ErrorLayerNotPresent:
		errName = "ErrorLayerNotPresent"
	case vulkan.ErrorExtensionNotPresent:
		errName = "ErrorExtensionNotPresent"
	case vulkan.ErrorFeatureNotPresent:
		errName = "ErrorFeatureNotPresent"
	case vulkan.ErrorIncompatibleDriver:
		errName = "ErrorIncompatibleDriver"
	case vulkan.ErrorTooManyObjects:
		errName = "ErrorTooManyObjects"
	case vulkan.ErrorFormatNotSupported:
		errName = "ErrorFormatNotSupported"
	case vulkan.ErrorFragmentedPool:
		errName = "ErrorFragmentedPool"
	case vulkan.ErrorOutOfPoolMemory:
		errName = "ErrorOutOfPoolMemory"
	case vulkan.ErrorInvalidExternalHandle:
		errName = "ErrorInvalidExternalHandle"
	case vulkan.ErrorSurfaceLost:
		errName = "ErrorSurfaceLost"
	case vulkan.ErrorNativeWindowInUse:
		errName = "ErrorNativeWindowInUse"
	case vulkan.Suboptimal:
		errName = "Suboptimal"
	case vulkan.ErrorOutOfDate:
		errName = "ErrorOutOfDate"
	case vulkan.ErrorIncompatibleDisplay:
		errName = "ErrorIncompatibleDisplay"
	case vulkan.ErrorValidationFailed:
		errName = "ErrorValidationFailed"
	case vulkan.ErrorInvalidShaderNv:
		errName = "ErrorInvalidShaderNv"
	case vulkan.ErrorInvalidDrmFormatModifierPlaneLayout:
		errName = "ErrorInvalidDrmFormatModifierPlaneLayout"
	case vulkan.ErrorFragmentation:
		errName = "ErrorFragmentation"
	case vulkan.ErrorNotPermitted:
		errName = "ErrorNotPermitted"
	default:
		errName = "unknown code"
	}

	return fmt.Sprintf("%s (%d)", errName, error)
}

func fmtObjectType(objectType vulkan.DebugReportObjectType) string {
	switch objectType {
	case vulkan.DebugReportObjectTypeInstance:
		return "Instance"
	case vulkan.DebugReportObjectTypePhysicalDevice:
		return "PhysicalDevice"
	case vulkan.DebugReportObjectTypeDevice:
		return "Device"
	case vulkan.DebugReportObjectTypeQueue:
		return "Queue"
	case vulkan.DebugReportObjectTypeSemaphore:
		return "Semaphore"
	case vulkan.DebugReportObjectTypeCommandBuffer:
		return "CommandBuffer"
	case vulkan.DebugReportObjectTypeFence:
		return "Fence"
	case vulkan.DebugReportObjectTypeDeviceMemory:
		return "DeviceMemory"
	case vulkan.DebugReportObjectTypeBuffer:
		return "Buffer"
	case vulkan.DebugReportObjectTypeImage:
		return "Image"
	case vulkan.DebugReportObjectTypeEvent:
		return "Event"
	case vulkan.DebugReportObjectTypeQueryPool:
		return "QueryPool"
	case vulkan.DebugReportObjectTypeBufferView:
		return "BufferView"
	case vulkan.DebugReportObjectTypeImageView:
		return "ImageView"
	case vulkan.DebugReportObjectTypeShaderModule:
		return "ShaderModule"
	case vulkan.DebugReportObjectTypePipelineCache:
		return "PipelineCache"
	case vulkan.DebugReportObjectTypePipelineLayout:
		return "PipelineLayout"
	case vulkan.DebugReportObjectTypeRenderPass:
		return "RenderPass"
	case vulkan.DebugReportObjectTypePipeline:
		return "Pipeline"
	case vulkan.DebugReportObjectTypeDescriptorSetLayout:
		return "DescriptorSetLayout"
	case vulkan.DebugReportObjectTypeSampler:
		return "Sampler"
	case vulkan.DebugReportObjectTypeDescriptorPool:
		return "DescriptorPool"
	case vulkan.DebugReportObjectTypeDescriptorSet:
		return "DescriptorSet"
	case vulkan.DebugReportObjectTypeFramebuffer:
		return "Framebuffer"
	case vulkan.DebugReportObjectTypeCommandPool:
		return "CommandPool"
	case vulkan.DebugReportObjectTypeSurfaceKhr:
		return "SurfaceKhr"
	case vulkan.DebugReportObjectTypeSwapchainKhr:
		return "SwapchainKhr"
	case vulkan.DebugReportObjectTypeDebugReportCallbackExt:
		return "DebugReport/DebugReportCallbackExt"
	case vulkan.DebugReportObjectTypeDisplayKhr:
		return "DisplayKhr"
	case vulkan.DebugReportObjectTypeDisplayModeKhr:
		return "DisplayModeKhr"
	case vulkan.DebugReportObjectTypeObjectTableNvx:
		return "ObjectTableNvx"
	case vulkan.DebugReportObjectTypeIndirectCommandsLayoutNvx:
		return "IndirectCommandsLayoutNvx"
	case vulkan.DebugReportObjectTypeValidationCacheExt:
		return "ValidationCache/ValidationCacheExt"
	case vulkan.DebugReportObjectTypeSamplerYcbcrConversion:
		return "SamplerYcbcrConversion/SamplerYcbcrConversionKhr"
	case vulkan.DebugReportObjectTypeDescriptorUpdateTemplate:
		return "DescriptorUpdateTemplate/DescriptorUpdateTemplateKhr"
	case vulkan.DebugReportObjectTypeAccelerationStructureNvx:
		return "AccelerationStructureNvx"
	default:
		return "Unknown"
	}
}
