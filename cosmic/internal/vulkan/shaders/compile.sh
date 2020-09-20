glslc shader.vert -o vert.spv
glslc shader.frag -o frag.spv

go-bindata -nocompress -pkg=shaders frag.spv vert.spv