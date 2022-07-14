// Code generated for package bindata by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/templates/cluster_op.html
// pkg/templates/cluster_op.md
// pkg/templates/test.html
package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

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

var _pkgTemplatesCluster_opHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x94\x41\x6f\xd3\x30\x14\xc7\xef\xfd\x14\x0f\xef\xdc\xa5\x2b\x20\x41\x70\x2a\x55\x94\x5d\x90\xe8\x34\x15\x4d\x1c\x9d\xf8\xb5\x35\xb8\x71\xe4\xb8\x5b\x4b\x94\x13\x62\x42\x42\x9c\x38\x4c\x53\x27\x21\x24\xe0\xc8\x2e\x08\x04\x13\x7c\x1a\xd6\xae\xdf\x02\xd9\x49\xb7\x76\x27\x0e\xab\x2d\xc5\x79\xf9\xbf\xfc\x7f\xcf\x8e\x5e\xe8\xad\x56\xfb\x61\xe7\xd9\xce\x23\xe8\x9b\x81\x6c\x54\xa8\x5d\x40\xb2\xb8\x17\x10\x8c\x89\x7d\x80\x8c\x37\x2a\x00\x00\x74\x80\x86\x41\xd4\x67\x3a\x45\x13\x90\xa7\x9d\xed\xea\x3d\x52\x4a\x46\x18\x89\x8d\x8e\xbd\x52\xaf\x08\x2a\xd4\x2b\xde\xa5\xa9\x19\xdb\xd8\x26\x1a\x16\x4a\x84\xcc\xdd\xdb\xd1\x55\xb1\xa9\xa6\xe2\x25\xfa\xb0\x75\x27\x19\x3d\xb8\x14\x5c\x62\x55\xb2\xb1\x1a\x1a\x1f\xba\x62\x84\xfc\x4a\x0c\x95\xe6\xa8\x7d\xb8\x9b\x8c\x20\x55\x52\x70\xd8\xd8\xae\xdb\x59\xa4\xe4\x95\x82\xc5\xaf\x83\xba\x6c\x20\xe4\xd8\x87\xa6\x16\x4c\x5e\xd9\xed\xb5\x77\x5b\xd5\xbd\xdd\xe6\x8e\x0f\xa1\x46\xf6\xa2\x7a\xa0\x34\x5f\x58\x51\xaf\x2c\x9f\x86\x8a\x8f\x1b\x15\xca\xc5\x3e\x30\x29\x7a\x71\x40\x22\x8c\x0d\xea\xcb\x33\x70\x7b\x2b\x6a\x0b\x48\x8d\x40\x84\x52\xa6\x09\x8b\x84\x3d\xcd\x7a\x11\x27\x8c\xf3\x45\x7c\x20\xb8\xe9\x07\xe4\x7e\xad\x56\x5a\x14\x36\x1a\xc2\x5e\xa4\xa4\xd2\x01\xd9\x68\x6d\xd9\xb9\x24\x17\x29\xfd\x45\x05\x12\xbb\x86\x80\xab\x30\x20\x4b\x87\x59\xbf\x9d\x8c\x48\x23\xcb\x36\xdd\xb7\xc8\x73\xea\x99\xfe\x12\xc3\x33\x7a\x85\x78\xdd\x9f\xaf\xf8\xaf\xaa\x76\xcc\x3f\xfe\x9c\x4d\xbe\xfa\x90\x65\x9b\x89\x56\xcf\x31\x32\x4f\xd8\x00\xf3\x7c\xd5\xc6\x33\xfc\x26\x99\x93\xc3\xd9\x9f\x4f\x8e\xa9\x31\x55\x43\x1d\xe1\xfa\xa1\xd3\xf7\xef\xfe\xfe\x3e\x71\x50\x95\xa0\x66\x46\xe9\xf5\x02\xcf\xdf\x4c\xce\xcf\x7e\x4d\x8f\xbe\xcf\x8f\xbe\x39\x6c\xa4\x91\x19\xe4\x4d\xb3\x5e\xee\xc5\xe9\x8f\xd9\xf1\xd9\xfc\xf8\x10\x1e\x0f\x43\x6c\x97\x7b\x85\xe9\x87\xcf\xb3\x93\xb7\x17\xa7\x5f\xa6\xaf\x5e\xff\x17\x9e\x7a\xae\x0d\x6c\xf7\x73\xb1\x6f\x97\xb2\x6d\x3c\xf7\x83\xf9\x17\x00\x00\xff\xff\x5e\x6c\x05\xb6\x70\x04\x00\x00")

func pkgTemplatesCluster_opHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pkgTemplatesCluster_opHtml,
		"pkg/templates/cluster_op.html",
	)
}

func pkgTemplatesCluster_opHtml() (*asset, error) {
	bytes, err := pkgTemplatesCluster_opHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/templates/cluster_op.html", size: 1136, mode: os.FileMode(420), modTime: time.Unix(1657782569, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pkgTemplatesCluster_opMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x56\x56\xa8\xae\xd6\x2b\xc9\x2c\xc9\x49\xad\xad\xe5\xd2\xd2\x7a\xb9\x70\xe7\xf3\xd9\xeb\xb4\xb4\x14\xaa\xab\x15\xf4\x0a\x8a\xf2\xb3\x52\x93\x4b\xfc\x12\x73\x53\x15\x20\xb2\xb3\xdb\x9e\xef\x5b\x02\x95\x2d\x4a\x2d\xce\x2f\x2d\x4a\x4e\x45\x48\x3f\x9b\xbe\xed\xe5\xf4\x2d\x50\xe9\xe4\xa2\xd4\xc4\x92\xd4\x14\xc7\x12\x88\xd4\xe4\xde\x27\x7b\xe7\x68\x69\x59\x81\xac\xcb\x2f\x48\x2d\x4a\x2c\xc9\x2f\xaa\xad\xe5\xb2\x49\xcb\xcf\x2b\x51\x48\xce\xcf\xc9\x2f\xb2\x55\xca\xcc\x4b\xcb\x57\xb2\x7b\x36\x67\xcd\xb3\x6d\x1d\xcf\x1a\xd7\x3f\x9f\xb2\xd1\xbb\x34\x29\xd5\x1f\xaa\xfa\x45\xfb\xaa\xa7\x5d\x2b\x9e\xf6\x4f\x7c\xd9\xd0\x68\xa3\x0f\xd2\x67\x07\x08\x00\x00\xff\xff\x21\xbf\x9a\x53\xbf\x00\x00\x00")

func pkgTemplatesCluster_opMdBytes() ([]byte, error) {
	return bindataRead(
		_pkgTemplatesCluster_opMd,
		"pkg/templates/cluster_op.md",
	)
}

func pkgTemplatesCluster_opMd() (*asset, error) {
	bytes, err := pkgTemplatesCluster_opMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/templates/cluster_op.md", size: 191, mode: os.FileMode(420), modTime: time.Unix(1657783128, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pkgTemplatesTestHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x56\xd0\xcb\x4d\x2d\x2e\x4e\x4c\x4f\x55\xa8\xad\x55\xd0\x55\x00\x09\xa4\x24\x96\x80\x78\x80\x00\x00\x00\xff\xff\xd4\x42\x3a\xbf\x1c\x00\x00\x00")

func pkgTemplatesTestHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pkgTemplatesTestHtml,
		"pkg/templates/test.html",
	)
}

func pkgTemplatesTestHtml() (*asset, error) {
	bytes, err := pkgTemplatesTestHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/templates/test.html", size: 28, mode: os.FileMode(420), modTime: time.Unix(1657705532, 0)}
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
	"pkg/templates/cluster_op.html": pkgTemplatesCluster_opHtml,
	"pkg/templates/cluster_op.md":   pkgTemplatesCluster_opMd,
	"pkg/templates/test.html":       pkgTemplatesTestHtml,
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
	"pkg": &bintree{nil, map[string]*bintree{
		"templates": &bintree{nil, map[string]*bintree{
			"cluster_op.html": &bintree{pkgTemplatesCluster_opHtml, map[string]*bintree{}},
			"cluster_op.md":   &bintree{pkgTemplatesCluster_opMd, map[string]*bintree{}},
			"test.html":       &bintree{pkgTemplatesTestHtml, map[string]*bintree{}},
		}},
	}},
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
