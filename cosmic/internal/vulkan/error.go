package vulkan

import (
	"fmt"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

func panicOnError(result vulkan.Result, operation string) {
	if result < vulkan.Success {
		log.PanicfCore("failed to %s: %s", operation, operation, fmtResult(result))
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

	return fmt.Sprintf("%d (%s)", error, errName)
}
