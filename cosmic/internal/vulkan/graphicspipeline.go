package vulkan

import (
	"github.com/lentus/cosmic-engine/cosmic/log"
	"github.com/vulkan-go/vulkan"
)

func (ctx *Context) createGraphicsPipeline() {
	vertexShaderModule := ctx.createShaderModule("vert.spv")
	fragmentShaderModule := ctx.createShaderModule("frag.spv")

	vertexShaderStageCreateInfo := vulkan.PipelineShaderStageCreateInfo{
		SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
		Stage:  vulkan.ShaderStageVertexBit,
		Module: vertexShaderModule,
		PName:  safeStr("main"),
	}

	fragmentShaderStageCreateInfo := vulkan.PipelineShaderStageCreateInfo{
		SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
		Stage:  vulkan.ShaderStageFragmentBit,
		Module: fragmentShaderModule,
		PName:  safeStr("main"),
	}

	vertexInputStateCreateInfo := vulkan.PipelineVertexInputStateCreateInfo{
		SType:                           vulkan.StructureTypePipelineVertexInputStateCreateInfo,
		VertexBindingDescriptionCount:   0,
		PVertexBindingDescriptions:      nil,
		VertexAttributeDescriptionCount: 0,
		PVertexAttributeDescriptions:    nil,
	}

	inputAssemblyStateCreateInfo := vulkan.PipelineInputAssemblyStateCreateInfo{
		SType:                  vulkan.StructureTypePipelineInputAssemblyStateCreateInfo,
		Topology:               vulkan.PrimitiveTopologyTriangleList,
		PrimitiveRestartEnable: vulkan.False,
	}

	viewport := vulkan.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(ctx.swapchainImageExtent.Width),
		Height:   float32(ctx.swapchainImageExtent.Height),
		MinDepth: 0,
		MaxDepth: 1,
	}

	scissor := vulkan.Rect2D{
		Offset: vulkan.Offset2D{X: 0, Y: 0},
		Extent: vulkan.Extent2D{
			Width:  ctx.swapchainImageExtent.Width,
			Height: ctx.swapchainImageExtent.Height,
		},
	}

	viewportStateCreateInfo := vulkan.PipelineViewportStateCreateInfo{
		SType:         vulkan.StructureTypePipelineViewportStateCreateInfo,
		ViewportCount: 1,
		PViewports:    []vulkan.Viewport{viewport},
		ScissorCount:  1,
		PScissors:     []vulkan.Rect2D{scissor},
	}

	rasterizationStateCreateInfo := vulkan.PipelineRasterizationStateCreateInfo{
		SType:                   vulkan.StructureTypePipelineRasterizationStateCreateInfo,
		DepthClampEnable:        vulkan.False,
		RasterizerDiscardEnable: vulkan.False,
		PolygonMode:             vulkan.PolygonModeFill,
		LineWidth:               1,
		CullMode:                vulkan.CullModeFlags(vulkan.CullModeBackBit),
		FrontFace:               vulkan.FrontFaceClockwise,
		DepthBiasEnable:         vulkan.False,
		DepthBiasConstantFactor: 0,
		DepthBiasClamp:          0,
		DepthBiasSlopeFactor:    0,
	}

	multisampleStateCreateInfo := vulkan.PipelineMultisampleStateCreateInfo{
		SType:                 vulkan.StructureTypePipelineMultisampleStateCreateInfo,
		SampleShadingEnable:   vulkan.False,
		RasterizationSamples:  vulkan.SampleCount1Bit,
		MinSampleShading:      1,
		PSampleMask:           nil,
		AlphaToCoverageEnable: vulkan.False,
		AlphaToOneEnable:      vulkan.False,
	}

	//ctx.createDepthStencilImage()

	colorComponentFlags := vulkan.ColorComponentFlags(
		vulkan.ColorComponentRBit | vulkan.ColorComponentGBit | vulkan.ColorComponentBBit | vulkan.ColorComponentABit,
	)
	colorblendAttachment := vulkan.PipelineColorBlendAttachmentState{
		BlendEnable:         vulkan.False,
		SrcColorBlendFactor: vulkan.BlendFactorOne,
		DstColorBlendFactor: vulkan.BlendFactorZero,
		ColorBlendOp:        vulkan.BlendOpAdd,
		SrcAlphaBlendFactor: vulkan.BlendFactorOne,
		DstAlphaBlendFactor: vulkan.BlendFactorZero,
		AlphaBlendOp:        vulkan.BlendOpAdd,
		ColorWriteMask:      colorComponentFlags,
	}

	colorblendCreateInfo := vulkan.PipelineColorBlendStateCreateInfo{
		SType:           vulkan.StructureTypePipelineColorBlendStateCreateInfo,
		LogicOpEnable:   vulkan.False,
		LogicOp:         vulkan.LogicOpCopy,
		AttachmentCount: 1,
		PAttachments:    []vulkan.PipelineColorBlendAttachmentState{colorblendAttachment},
		BlendConstants:  [4]float32{0, 0, 0, 0},
	}

	pipelineLayoutCreateInfo := vulkan.PipelineLayoutCreateInfo{
		SType:                  vulkan.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount:         0,
		PSetLayouts:            nil,
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}
	var pipelineLayout vulkan.PipelineLayout
	result := vulkan.CreatePipelineLayout(ctx.device, &pipelineLayoutCreateInfo, nil, &pipelineLayout)
	panicOnError(result, "create pipeline layout")
	ctx.pipelineLayout = pipelineLayout

	pipelineCreateInfo := vulkan.GraphicsPipelineCreateInfo{
		SType:      vulkan.StructureTypeGraphicsPipelineCreateInfo,
		StageCount: 2,
		PStages: []vulkan.PipelineShaderStageCreateInfo{
			vertexShaderStageCreateInfo, fragmentShaderStageCreateInfo,
		},
		PVertexInputState:   &vertexInputStateCreateInfo,
		PInputAssemblyState: &inputAssemblyStateCreateInfo,
		PTessellationState:  nil,
		PViewportState:      &viewportStateCreateInfo,
		PRasterizationState: &rasterizationStateCreateInfo,
		PMultisampleState:   &multisampleStateCreateInfo,
		PDepthStencilState:  nil,
		PColorBlendState:    &colorblendCreateInfo,
		PDynamicState:       nil,
		Layout:              pipelineLayout,
		RenderPass:          ctx.renderPass,
		Subpass:             0,
		BasePipelineHandle:  nil,
		BasePipelineIndex:   -1,
	}

	graphicsPipelines := make([]vulkan.Pipeline, 1)
	pipelineCreateInfos := []vulkan.GraphicsPipelineCreateInfo{pipelineCreateInfo}
	result = vulkan.CreateGraphicsPipelines(
		ctx.device, vulkan.NullPipelineCache, 1, pipelineCreateInfos, nil, graphicsPipelines,
	)
	panicOnError(result, "create graphics pipeline")
	ctx.graphicsPipeline = graphicsPipelines[0]

	vulkan.DestroyShaderModule(ctx.device, vertexShaderModule, nil)
	vulkan.DestroyShaderModule(ctx.device, fragmentShaderModule, nil)
}

func (ctx *Context) destroyGraphicsPipeline() {
	vulkan.DestroyPipeline(ctx.device, ctx.graphicsPipeline, nil)
	vulkan.DestroyPipelineLayout(ctx.device, ctx.pipelineLayout, nil)
	//ctx.destroyDepthStencilImage()
	vulkan.DestroyRenderPass(ctx.device, ctx.renderPass, nil)
}

func (ctx *Context) createDepthStencilImage() {
	log.DebugCore("Creating Vulkan depth stencil image")

	// Take the first supported format of the following formats
	desiredFormats := []vulkan.Format{
		vulkan.FormatD32SfloatS8Uint,
		vulkan.FormatD24UnormS8Uint,
		vulkan.FormatD16UnormS8Uint,
		vulkan.FormatD32Sfloat,
		vulkan.FormatD16Unorm,
	}
	for _, format := range desiredFormats {
		var formatProps vulkan.FormatProperties
		vulkan.GetPhysicalDeviceFormatProperties(ctx.gpu.ref, format, &formatProps)
		formatProps.Deref()

		if formatProps.OptimalTilingFeatures&vulkan.FormatFeatureFlags(vulkan.FormatFeatureDepthStencilAttachmentBit) != 0 {
			ctx.depthStencilFormat = format
			break
		}
	}

	if ctx.depthStencilFormat == vulkan.FormatUndefined {
		log.PanicCore("None of the desired depth stencil formats are supported")
	}

	// Check whether stencil is available
	ctx.stencilAvailable = ctx.depthStencilFormat == vulkan.FormatD32SfloatS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD24UnormS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD16UnormS8Uint ||
		ctx.depthStencilFormat == vulkan.FormatD32Sfloat

	imageCreateInfo := vulkan.ImageCreateInfo{
		SType:     vulkan.StructureTypeImageCreateInfo,
		Flags:     0,
		ImageType: vulkan.ImageType2d,
		Format:    ctx.depthStencilFormat,
		Extent: vulkan.Extent3D{
			Width:  ctx.surface.capabilities.CurrentExtent.Width,
			Height: ctx.surface.capabilities.CurrentExtent.Height,
			Depth:  1,
		},
		MipLevels:             1,
		ArrayLayers:           1,
		Samples:               vulkan.SampleCount1Bit,
		Tiling:                vulkan.ImageTilingOptimal,
		Usage:                 vulkan.ImageUsageFlags(vulkan.ImageUsageDepthStencilAttachmentBit),
		SharingMode:           vulkan.SharingModeExclusive,
		QueueFamilyIndexCount: 0,   // Ignored because of exclusive mode
		PQueueFamilyIndices:   nil, // Ignored because of exclusive mode
		InitialLayout:         vulkan.ImageLayoutUndefined,
	}

	var depthStencilImage vulkan.Image
	result := vulkan.CreateImage(ctx.device, &imageCreateInfo, nil, &depthStencilImage)
	panicOnError(result, "create depth stencil image")
	ctx.depthStencilImage = depthStencilImage

	var imageMemoryRequirements vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(ctx.device, ctx.depthStencilImage, &imageMemoryRequirements)
	imageMemoryRequirements.Deref()

	memoryTypeIndex := ctx.findMemoryTypeIndex(&imageMemoryRequirements, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit))
	if memoryTypeIndex == vulkan.MaxUint32 {
		log.PanicCore("Could not find memory type to allocate depth stencil image memory")
	}

	memoryAllocateInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  imageMemoryRequirements.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}
	var depthStencilImageMemory vulkan.DeviceMemory
	vulkan.AllocateMemory(ctx.device, &memoryAllocateInfo, nil, &depthStencilImageMemory)
	ctx.depthStencilImageMemory = depthStencilImageMemory
	vulkan.BindImageMemory(ctx.device, ctx.depthStencilImage, ctx.depthStencilImageMemory, 0)

	aspectMask := vulkan.ImageAspectDepthBit
	if ctx.stencilAvailable {
		aspectMask |= vulkan.ImageAspectStencilBit
	}

	imageViewCreateInfo := vulkan.ImageViewCreateInfo{
		SType:      vulkan.StructureTypeImageViewCreateInfo,
		Image:      ctx.depthStencilImage,
		ViewType:   vulkan.ImageViewType2d,
		Format:     ctx.depthStencilFormat,
		Components: vulkan.ComponentMapping{}, // Use identity mapping for rgba components
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectFlags(aspectMask),
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var depthStencilImageView vulkan.ImageView
	result = vulkan.CreateImageView(ctx.device, &imageViewCreateInfo, nil, &depthStencilImageView)
	panicOnError(result, "create depth stencil image view")
	ctx.depthStencilImageView = depthStencilImageView
}

func (ctx *Context) destroyDepthStencilImage() {
	vulkan.DestroyImageView(ctx.device, ctx.depthStencilImageView, nil)
	vulkan.FreeMemory(ctx.device, ctx.depthStencilImageMemory, nil)
	vulkan.DestroyImage(ctx.device, ctx.depthStencilImage, nil)
}
