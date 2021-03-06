// Code generated for package shaders by go-bindata DO NOT EDIT. (@generated)
// sources:
// frag.spv
// vert.spv
package shaders

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _fragSpv = []byte("\x03\x02#\a\x00\x00\x01\x00\n\x00\r\x00\x13\x00\x00\x00\x00\x00\x00\x00\x11\x00\x02\x00\x01\x00\x00\x00\v\x00\x06\x00\x01\x00\x00\x00GLSL.std.450\x00\x00\x00\x00\x0e\x00\x03\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0f\x00\a\x00\x04\x00\x00\x00\x04\x00\x00\x00main\x00\x00\x00\x00\t\x00\x00\x00\f\x00\x00\x00\x10\x00\x03\x00\x04\x00\x00\x00\a\x00\x00\x00\x03\x00\x03\x00\x02\x00\x00\x00\xc2\x01\x00\x00\x04\x00\t\x00GL_ARB_separate_shader_objects\x00\x00\x04\x00\n\x00GL_GOOGLE_cpp_style_line_directive\x00\x00\x04\x00\b\x00GL_GOOGLE_include_directive\x00\x05\x00\x04\x00\x04\x00\x00\x00main\x00\x00\x00\x00\x05\x00\x05\x00\t\x00\x00\x00outColor\x00\x00\x00\x00\x05\x00\x05\x00\f\x00\x00\x00fragColor\x00\x00\x00G\x00\x04\x00\t\x00\x00\x00\x1e\x00\x00\x00\x00\x00\x00\x00G\x00\x04\x00\f\x00\x00\x00\x1e\x00\x00\x00\x00\x00\x00\x00\x13\x00\x02\x00\x02\x00\x00\x00!\x00\x03\x00\x03\x00\x00\x00\x02\x00\x00\x00\x16\x00\x03\x00\x06\x00\x00\x00 \x00\x00\x00\x17\x00\x04\x00\a\x00\x00\x00\x06\x00\x00\x00\x04\x00\x00\x00 \x00\x04\x00\b\x00\x00\x00\x03\x00\x00\x00\a\x00\x00\x00;\x00\x04\x00\b\x00\x00\x00\t\x00\x00\x00\x03\x00\x00\x00\x17\x00\x04\x00\n\x00\x00\x00\x06\x00\x00\x00\x03\x00\x00\x00 \x00\x04\x00\v\x00\x00\x00\x01\x00\x00\x00\n\x00\x00\x00;\x00\x04\x00\v\x00\x00\x00\f\x00\x00\x00\x01\x00\x00\x00+\x00\x04\x00\x06\x00\x00\x00\x0e\x00\x00\x00\x00\x00\x80?6\x00\x05\x00\x02\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\xf8\x00\x02\x00\x05\x00\x00\x00=\x00\x04\x00\n\x00\x00\x00\r\x00\x00\x00\f\x00\x00\x00Q\x00\x05\x00\x06\x00\x00\x00\x0f\x00\x00\x00\r\x00\x00\x00\x00\x00\x00\x00Q\x00\x05\x00\x06\x00\x00\x00\x10\x00\x00\x00\r\x00\x00\x00\x01\x00\x00\x00Q\x00\x05\x00\x06\x00\x00\x00\x11\x00\x00\x00\r\x00\x00\x00\x02\x00\x00\x00P\x00\a\x00\a\x00\x00\x00\x12\x00\x00\x00\x0f\x00\x00\x00\x10\x00\x00\x00\x11\x00\x00\x00\x0e\x00\x00\x00>\x00\x03\x00\t\x00\x00\x00\x12\x00\x00\x00\xfd\x00\x01\x008\x00\x01\x00")

func fragSpvBytes() ([]byte, error) {
	return _fragSpv, nil
}

func fragSpv() (*asset, error) {
	bytes, err := fragSpvBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "frag.spv", size: 608, mode: os.FileMode(436), modTime: time.Unix(1600449106, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vertSpv = []byte("\x03\x02#\a\x00\x00\x01\x00\n\x00\r\x006\x00\x00\x00\x00\x00\x00\x00\x11\x00\x02\x00\x01\x00\x00\x00\v\x00\x06\x00\x01\x00\x00\x00GLSL.std.450\x00\x00\x00\x00\x0e\x00\x03\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0f\x00\b\x00\x00\x00\x00\x00\x04\x00\x00\x00main\x00\x00\x00\x00\"\x00\x00\x00&\x00\x00\x001\x00\x00\x00\x03\x00\x03\x00\x02\x00\x00\x00\xc2\x01\x00\x00\x04\x00\t\x00GL_ARB_separate_shader_objects\x00\x00\x04\x00\n\x00GL_GOOGLE_cpp_style_line_directive\x00\x00\x04\x00\b\x00GL_GOOGLE_include_directive\x00\x05\x00\x04\x00\x04\x00\x00\x00main\x00\x00\x00\x00\x05\x00\x05\x00\f\x00\x00\x00positions\x00\x00\x00\x05\x00\x04\x00\x17\x00\x00\x00colors\x00\x00\x05\x00\x06\x00 \x00\x00\x00gl_PerVertex\x00\x00\x00\x00\x06\x00\x06\x00 \x00\x00\x00\x00\x00\x00\x00gl_Position\x00\x06\x00\a\x00 \x00\x00\x00\x01\x00\x00\x00gl_PointSize\x00\x00\x00\x00\x06\x00\a\x00 \x00\x00\x00\x02\x00\x00\x00gl_ClipDistance\x00\x06\x00\a\x00 \x00\x00\x00\x03\x00\x00\x00gl_CullDistance\x00\x05\x00\x03\x00\"\x00\x00\x00\x00\x00\x00\x00\x05\x00\x06\x00&\x00\x00\x00gl_VertexIndex\x00\x00\x05\x00\x05\x001\x00\x00\x00fragColor\x00\x00\x00H\x00\x05\x00 \x00\x00\x00\x00\x00\x00\x00\v\x00\x00\x00\x00\x00\x00\x00H\x00\x05\x00 \x00\x00\x00\x01\x00\x00\x00\v\x00\x00\x00\x01\x00\x00\x00H\x00\x05\x00 \x00\x00\x00\x02\x00\x00\x00\v\x00\x00\x00\x03\x00\x00\x00H\x00\x05\x00 \x00\x00\x00\x03\x00\x00\x00\v\x00\x00\x00\x04\x00\x00\x00G\x00\x03\x00 \x00\x00\x00\x02\x00\x00\x00G\x00\x04\x00&\x00\x00\x00\v\x00\x00\x00*\x00\x00\x00G\x00\x04\x001\x00\x00\x00\x1e\x00\x00\x00\x00\x00\x00\x00\x13\x00\x02\x00\x02\x00\x00\x00!\x00\x03\x00\x03\x00\x00\x00\x02\x00\x00\x00\x16\x00\x03\x00\x06\x00\x00\x00 \x00\x00\x00\x17\x00\x04\x00\a\x00\x00\x00\x06\x00\x00\x00\x02\x00\x00\x00\x15\x00\x04\x00\b\x00\x00\x00 \x00\x00\x00\x00\x00\x00\x00+\x00\x04\x00\b\x00\x00\x00\t\x00\x00\x00\x03\x00\x00\x00\x1c\x00\x04\x00\n\x00\x00\x00\a\x00\x00\x00\t\x00\x00\x00 \x00\x04\x00\v\x00\x00\x00\x06\x00\x00\x00\n\x00\x00\x00;\x00\x04\x00\v\x00\x00\x00\f\x00\x00\x00\x06\x00\x00\x00+\x00\x04\x00\x06\x00\x00\x00\r\x00\x00\x00\x00\x00\x00\x00+\x00\x04\x00\x06\x00\x00\x00\x0e\x00\x00\x00\x00\x00\x00\xbf,\x00\x05\x00\a\x00\x00\x00\x0f\x00\x00\x00\r\x00\x00\x00\x0e\x00\x00\x00+\x00\x04\x00\x06\x00\x00\x00\x10\x00\x00\x00\x00\x00\x00?,\x00\x05\x00\a\x00\x00\x00\x11\x00\x00\x00\x10\x00\x00\x00\x10\x00\x00\x00,\x00\x05\x00\a\x00\x00\x00\x12\x00\x00\x00\x0e\x00\x00\x00\x10\x00\x00\x00,\x00\x06\x00\n\x00\x00\x00\x13\x00\x00\x00\x0f\x00\x00\x00\x11\x00\x00\x00\x12\x00\x00\x00\x17\x00\x04\x00\x14\x00\x00\x00\x06\x00\x00\x00\x03\x00\x00\x00\x1c\x00\x04\x00\x15\x00\x00\x00\x14\x00\x00\x00\t\x00\x00\x00 \x00\x04\x00\x16\x00\x00\x00\x06\x00\x00\x00\x15\x00\x00\x00;\x00\x04\x00\x16\x00\x00\x00\x17\x00\x00\x00\x06\x00\x00\x00+\x00\x04\x00\x06\x00\x00\x00\x18\x00\x00\x00\x00\x00\x80?,\x00\x06\x00\x14\x00\x00\x00\x19\x00\x00\x00\x18\x00\x00\x00\r\x00\x00\x00\r\x00\x00\x00,\x00\x06\x00\x14\x00\x00\x00\x1a\x00\x00\x00\r\x00\x00\x00\x18\x00\x00\x00\r\x00\x00\x00,\x00\x06\x00\x14\x00\x00\x00\x1b\x00\x00\x00\r\x00\x00\x00\r\x00\x00\x00\x18\x00\x00\x00,\x00\x06\x00\x15\x00\x00\x00\x1c\x00\x00\x00\x19\x00\x00\x00\x1a\x00\x00\x00\x1b\x00\x00\x00\x17\x00\x04\x00\x1d\x00\x00\x00\x06\x00\x00\x00\x04\x00\x00\x00+\x00\x04\x00\b\x00\x00\x00\x1e\x00\x00\x00\x01\x00\x00\x00\x1c\x00\x04\x00\x1f\x00\x00\x00\x06\x00\x00\x00\x1e\x00\x00\x00\x1e\x00\x06\x00 \x00\x00\x00\x1d\x00\x00\x00\x06\x00\x00\x00\x1f\x00\x00\x00\x1f\x00\x00\x00 \x00\x04\x00!\x00\x00\x00\x03\x00\x00\x00 \x00\x00\x00;\x00\x04\x00!\x00\x00\x00\"\x00\x00\x00\x03\x00\x00\x00\x15\x00\x04\x00#\x00\x00\x00 \x00\x00\x00\x01\x00\x00\x00+\x00\x04\x00#\x00\x00\x00$\x00\x00\x00\x00\x00\x00\x00 \x00\x04\x00%\x00\x00\x00\x01\x00\x00\x00#\x00\x00\x00;\x00\x04\x00%\x00\x00\x00&\x00\x00\x00\x01\x00\x00\x00 \x00\x04\x00(\x00\x00\x00\x06\x00\x00\x00\a\x00\x00\x00 \x00\x04\x00.\x00\x00\x00\x03\x00\x00\x00\x1d\x00\x00\x00 \x00\x04\x000\x00\x00\x00\x03\x00\x00\x00\x14\x00\x00\x00;\x00\x04\x000\x00\x00\x001\x00\x00\x00\x03\x00\x00\x00 \x00\x04\x003\x00\x00\x00\x06\x00\x00\x00\x14\x00\x00\x006\x00\x05\x00\x02\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\xf8\x00\x02\x00\x05\x00\x00\x00>\x00\x03\x00\f\x00\x00\x00\x13\x00\x00\x00>\x00\x03\x00\x17\x00\x00\x00\x1c\x00\x00\x00=\x00\x04\x00#\x00\x00\x00'\x00\x00\x00&\x00\x00\x00A\x00\x05\x00(\x00\x00\x00)\x00\x00\x00\f\x00\x00\x00'\x00\x00\x00=\x00\x04\x00\a\x00\x00\x00*\x00\x00\x00)\x00\x00\x00Q\x00\x05\x00\x06\x00\x00\x00+\x00\x00\x00*\x00\x00\x00\x00\x00\x00\x00Q\x00\x05\x00\x06\x00\x00\x00,\x00\x00\x00*\x00\x00\x00\x01\x00\x00\x00P\x00\a\x00\x1d\x00\x00\x00-\x00\x00\x00+\x00\x00\x00,\x00\x00\x00\r\x00\x00\x00\x18\x00\x00\x00A\x00\x05\x00.\x00\x00\x00/\x00\x00\x00\"\x00\x00\x00$\x00\x00\x00>\x00\x03\x00/\x00\x00\x00-\x00\x00\x00=\x00\x04\x00#\x00\x00\x002\x00\x00\x00&\x00\x00\x00A\x00\x05\x003\x00\x00\x004\x00\x00\x00\x17\x00\x00\x002\x00\x00\x00=\x00\x04\x00\x14\x00\x00\x005\x00\x00\x004\x00\x00\x00>\x00\x03\x001\x00\x00\x005\x00\x00\x00\xfd\x00\x01\x008\x00\x01\x00")

func vertSpvBytes() ([]byte, error) {
	return _vertSpv, nil
}

func vertSpv() (*asset, error) {
	bytes, err := vertSpvBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vert.spv", size: 1540, mode: os.FileMode(436), modTime: time.Unix(1600449106, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"frag.spv": fragSpv,
	"vert.spv": vertSpv,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"frag.spv": &bintree{fragSpv, map[string]*bintree{}},
	"vert.spv": &bintree{vertSpv, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
