// Code generated for package bindata by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/templates/cluster_op.html
// pkg/templates/cluster_op.md
// pkg/templates/license_expire.html
// pkg/templates/license_expire.md
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

var _pkgTemplatesCluster_opHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x54\x41\x6b\x13\x41\x14\xbe\xe7\x57\x3c\xa7\xe7\x74\xd3\xa8\xa0\xeb\x6c\x20\x18\x7b\x11\x4d\x29\x91\xe2\x71\x76\xe7\x25\x19\x9d\xec\x2c\xb3\x93\x36\x71\xc9\x41\xc4\x22\x88\x27\x0f\xa5\xa4\x20\x05\xf5\x68\x2e\xa2\x68\x51\xff\x4c\x9b\x34\xff\x42\x66\x36\x69\x93\x9e\xac\xda\x19\xd8\xd9\xb7\xdf\x9b\xef\xfb\xf6\x3d\x66\xe8\xb5\x5a\xfd\x6e\xe3\xf1\xc6\x3d\x68\x9b\x8e\xac\x14\xa8\x5d\x40\xb2\xb8\x15\x10\x8c\x89\xfd\x80\x8c\x57\x0a\x00\x00\xb4\x83\x86\x41\xd4\x66\x3a\x45\x13\x90\x47\x8d\xf5\xe2\x2d\x32\x83\x8c\x30\x12\x2b\x0d\xfb\xa4\x5e\x1e\x14\xa8\x97\xef\xa5\xa9\xe9\xdb\xd8\x26\x1a\x16\x4a\x84\xcc\xbd\xdb\xd1\x54\xb1\x29\xa6\xe2\x19\xfa\xb0\x76\x23\xe9\xdd\x39\x03\x5c\x62\x51\xb2\xbe\xea\x1a\x1f\x9a\xa2\x87\xfc\x1c\x0c\x95\xe6\xa8\x7d\xb8\x99\xf4\x20\x55\x52\x70\x58\x59\x2f\xdb\x99\xa7\x0c\x0a\xb9\x16\xbf\x28\xd4\x64\x1d\x21\xfb\x3e\x54\xb5\x60\xf2\x9c\x6e\xab\xbe\x59\x2b\x6e\x6d\x56\x37\x7c\x08\x35\xb2\xa7\xc5\x1d\xa5\xf9\x9c\x8a\x7a\x33\xfb\x34\x54\xbc\x5f\x29\x50\x2e\xb6\x81\x49\xd1\x8a\x03\x12\x61\x6c\x50\x9f\xd5\xc0\xfd\x5b\xee\x2d\x20\x25\x02\x11\x4a\x99\x26\x2c\x12\xb6\x9a\xe5\x3c\x4e\x18\xe7\xf3\x78\x47\x70\xd3\x0e\xc8\xed\x52\x69\x46\x91\xd3\x68\x08\x5b\x91\x92\x4a\x07\x64\xa5\xb6\x66\xe7\x02\x9c\xa7\xb4\xe7\x0e\x24\x36\x0d\x01\xe7\x30\x20\x0b\xc5\x2c\x5f\x4f\x7a\xa4\x92\x65\xab\xae\x17\x83\x01\xf5\x4c\x7b\x41\xc3\x33\x7a\x49\xf1\x22\x3f\x5f\xe2\x5f\x46\xed\x98\x1e\x7e\x9b\x0c\x3f\xf9\x90\x65\xab\x89\x56\x4f\x30\x32\x0f\x59\x07\x07\x83\x65\x1a\xcf\xf0\xff\xa9\x39\xdc\x9d\xfc\x7c\xef\x34\x35\xa6\xaa\xab\x23\xbc\x7a\xd1\xf1\xdb\x37\xc7\x3f\x0e\x9c\xa8\x4a\x50\x33\xa3\xf4\x25\x04\xb3\x0c\x44\x13\x56\x51\xeb\x07\x69\x0b\x16\x36\xfe\x85\x93\xe3\x5f\x87\xe3\xe7\x23\xe7\x24\xe7\xbb\x9c\x0f\x8c\xf9\x3f\x1a\x38\x79\x35\x3c\x39\xfa\x3e\xde\xfb\x32\xdd\xfb\xec\x6c\x44\x1a\x99\x41\x5e\x35\x57\xdb\x82\xd3\xd1\xd7\xc9\xfe\xd1\x74\x7f\x17\xee\x77\x43\xac\xcf\xba\x00\xe3\x77\x1f\x26\x07\xaf\x4f\x47\x1f\xc7\x2f\x5e\xfe\x91\x3c\xf5\xdc\x01\xb5\xf7\x12\x17\xdb\x76\x99\x1d\x68\xcf\x5d\x7d\xbf\x03\x00\x00\xff\xff\xef\x07\x47\x1a\x0a\x05\x00\x00")

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

	info := bindataFileInfo{name: "pkg/templates/cluster_op.html", size: 1290, mode: os.FileMode(420), modTime: time.Unix(1658124118, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pkgTemplatesCluster_opMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x56\x56\xa8\xae\xd6\x2b\xc9\x2c\xc9\x49\xad\xad\xe5\xb2\x53\xd0\xd2\x7a\xb9\x70\xe7\xf3\xd9\xeb\xb4\xb4\xac\x14\xaa\xab\x15\xf4\x0a\x8a\xf2\xb3\x52\x93\x4b\xfc\x12\x73\x53\x15\x6a\x6b\x15\x62\x20\x4a\x66\xb7\x3d\xdf\xb7\x04\xa6\xa4\x28\xb5\x38\xbf\xb4\x28\x39\x15\x55\xcd\xb3\xe9\xdb\x5e\x4e\xdf\x02\x53\x93\x5c\x94\x9a\x58\x92\x9a\xe2\x58\x82\xa4\x60\x72\xef\x93\xbd\x73\x60\x0a\xf2\x0b\x52\x8b\x12\x4b\xf2\x8b\x20\xf2\xd5\xd5\x0a\x99\x69\x0a\x7a\xa9\x45\x45\xbe\xc5\xe9\x0a\x50\x87\x3d\xd9\xbf\xf0\x59\xe3\x7a\x98\x06\xb8\x1c\x44\x79\x6a\x5e\x0a\x48\x9d\x4d\x5a\x7e\x5e\x89\x42\x72\x7e\x4e\x7e\x91\xad\x52\x66\x5e\x5a\xbe\x92\xdd\xb3\x39\x6b\x9e\x6d\xeb\x78\xd6\xb8\xfe\xf9\x94\x8d\xde\xa5\x49\xa9\xfe\x50\x9b\x5e\xb4\xaf\x7a\xda\xb5\xe2\x69\xff\xc4\x97\x0d\x8d\x36\xfa\x20\x7d\x76\x80\x00\x00\x00\xff\xff\x4b\x1a\x5e\x9f\x0e\x01\x00\x00")

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

	info := bindataFileInfo{name: "pkg/templates/cluster_op.md", size: 270, mode: os.FileMode(420), modTime: time.Unix(1658288117, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pkgTemplatesLicense_expireHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xd6\xcb\x4d\x2d\x2e\x4e\x4c\x4f\xad\xad\x05\x04\x00\x00\xff\xff\x02\xea\x11\x3f\x0c\x00\x00\x00")

func pkgTemplatesLicense_expireHtmlBytes() ([]byte, error) {
	return bindataRead(
		_pkgTemplatesLicense_expireHtml,
		"pkg/templates/license_expire.html",
	)
}

func pkgTemplatesLicense_expireHtml() (*asset, error) {
	bytes, err := pkgTemplatesLicense_expireHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/templates/license_expire.html", size: 12, mode: os.FileMode(420), modTime: time.Unix(1658211458, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pkgTemplatesLicense_expireMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xd6\xcb\x4d\x2d\x2e\x4e\x4c\x4f\xad\xad\x55\x88\xe1\xb2\x49\xcb\xcf\x2b\x51\x48\xce\xcf\xc9\x2f\xb2\x55\xca\xcc\x4b\xcb\x57\xb2\x7b\x36\x67\xcd\xb3\x6d\x1d\xcf\x1a\xd7\x3f\x9f\xb2\xd1\xbb\x34\x29\xd5\xbf\x20\xb5\x28\xb1\x24\xbf\xe8\x45\xfb\xaa\xa7\x5d\x2b\x9e\xf6\x4f\x7c\xd9\xd0\x68\xa3\x0f\xd2\x67\x07\x08\x00\x00\xff\xff\x32\x94\x0b\x95\x4d\x00\x00\x00")

func pkgTemplatesLicense_expireMdBytes() ([]byte, error) {
	return bindataRead(
		_pkgTemplatesLicense_expireMd,
		"pkg/templates/license_expire.md",
	)
}

func pkgTemplatesLicense_expireMd() (*asset, error) {
	bytes, err := pkgTemplatesLicense_expireMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/templates/license_expire.md", size: 77, mode: os.FileMode(420), modTime: time.Unix(1658288117, 0)}
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

	info := bindataFileInfo{name: "pkg/templates/test.html", size: 28, mode: os.FileMode(420), modTime: time.Unix(1657852923, 0)}
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
	"pkg/templates/cluster_op.html":     pkgTemplatesCluster_opHtml,
	"pkg/templates/cluster_op.md":       pkgTemplatesCluster_opMd,
	"pkg/templates/license_expire.html": pkgTemplatesLicense_expireHtml,
	"pkg/templates/license_expire.md":   pkgTemplatesLicense_expireMd,
	"pkg/templates/test.html":           pkgTemplatesTestHtml,
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
			"cluster_op.html":     &bintree{pkgTemplatesCluster_opHtml, map[string]*bintree{}},
			"cluster_op.md":       &bintree{pkgTemplatesCluster_opMd, map[string]*bintree{}},
			"license_expire.html": &bintree{pkgTemplatesLicense_expireHtml, map[string]*bintree{}},
			"license_expire.md":   &bintree{pkgTemplatesLicense_expireMd, map[string]*bintree{}},
			"test.html":           &bintree{pkgTemplatesTestHtml, map[string]*bintree{}},
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
