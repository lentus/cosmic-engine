package vulkan

import (
	"bytes"
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

func (ctx *Context) setupLayersAndExtensions() {
	requiredExtensions := ctx.nativeWindow.GetRequiredInstanceExtensions()

	log.DebugCore("Extensions required by GLFW:")
	for _, extension := range requiredExtensions {
		log.DebugfCore("\t%s", extension)
	}

	ctx.enabledInstanceExtensions = append(ctx.enabledInstanceExtensions, requiredExtensions...)
}

func (ctx *Context) getInstanceExtensions() []vulkan.ExtensionProperties {
	var extensionCount uint32
	vulkan.EnumerateInstanceExtensionProperties(safeStr(""), &extensionCount, nil)
	extensionPropertiesList := make([]vulkan.ExtensionProperties, extensionCount)
	result := vulkan.EnumerateInstanceExtensionProperties(safeStr(""), &extensionCount, extensionPropertiesList)
	panicOnError(result, "retrieve instance extensions")

	log.DebugfCore("Instance extensions (%d):", len(extensionPropertiesList))
	for _, props := range extensionPropertiesList {
		props.Deref()
		log.DebugfCore("\t%s [v%d]", props.ExtensionName, props.SpecVersion)
	}

	return extensionPropertiesList
}

func (ctx *Context) enableLayerIfAvailable(layerName string) {
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
