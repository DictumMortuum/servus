// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// assets/date.js
// assets/styles.css
// html/index.html
package main

import (
	"github.com/elazarl/go-bindata-assetfs"
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

var _assetsDateJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x91\xbf\xce\xda\x30\x14\xc5\x77\x3f\xc5\x51\x96\x2f\xd1\x87\x50\x3a\x74\x4a\x33\x50\xd1\x8a\xa2\x86\x0e\x65\x43\x0c\xc6\xb9\x10\x8b\xd8\x46\xce\x35\x15\xaa\x78\xf7\xca\x09\x90\xb4\x43\xb7\xe3\xf3\xcb\xfd\x93\x73\x95\xb3\x1d\xc3\x38\xcb\xcd\x46\x1a\xea\x50\x62\x27\x80\x64\x2d\x6d\x90\xfe\x96\xcc\x90\x7c\xa5\x83\x7f\xea\x4a\x7a\xd5\x44\xb1\xb8\x78\xdd\x0e\x4e\x0f\xd6\xc1\x52\x32\xeb\x2b\x43\xdb\x3b\x8b\x70\x0a\x1d\x47\xf5\x93\x2e\x4c\xe6\x40\x3e\x3e\x7e\x28\x76\x0f\xb9\x71\xd7\x97\xbd\x24\x35\x68\xb1\x2f\x84\x38\x06\xab\x58\x3b\x8b\x8e\xa5\xe7\xad\x36\x94\x66\xf8\x2d\x80\xab\xf4\x60\x57\xcb\x1b\x4a\x58\xfa\x85\xa5\x64\x4a\xb3\xe2\x41\x1a\x94\x03\x9d\x9f\x88\x57\x2e\xf8\x6e\x64\x66\xca\x2a\x6d\x03\xd3\x84\x56\x7f\xd1\x98\xc6\xc8\xea\x29\x1b\x07\xc6\x86\xaa\x21\x75\xee\xd7\x33\xbd\x57\x3b\x15\x0c\x59\x8e\x5f\x7e\x69\x29\xca\xcf\xb7\x6f\x75\xfa\xa6\x5a\xa7\xce\x6f\xd9\x5c\x5b\x4b\x7e\xb5\xad\xbe\xa3\x44\x83\x77\x24\x48\xf0\x0e\xf3\xdf\xda\x5a\x32\xfd\x53\x3a\xde\x6b\x57\xed\x5f\x6d\xea\xe7\xc6\x8c\x12\x1d\xf5\xb9\xb9\xc0\xe9\x2b\xc4\x19\x3e\xe6\x79\x56\x88\xfb\x24\xe1\xf1\x17\xf4\x10\xb1\x3e\x22\xd5\xf8\x84\x0f\xf9\xf0\x06\x34\x4a\x24\x79\x9c\xa0\x05\x70\x8f\x53\x3c\x71\xf0\x16\xba\xef\x35\x39\x52\x21\xfe\x04\x00\x00\xff\xff\x47\x85\xda\xa0\x51\x02\x00\x00")

func assetsDateJsBytes() ([]byte, error) {
	return bindataRead(
		_assetsDateJs,
		"assets/date.js",
	)
}

func assetsDateJs() (*asset, error) {
	bytes, err := assetsDateJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/date.js", size: 593, mode: os.FileMode(420), modTime: time.Unix(1575982167, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsStylesCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\xed\x8a\xdb\x3c\x10\x85\xff\xfb\x2a\x06\xcc\x0b\xef\x16\xbc\x58\xfe\x5a\x47\x7b\x35\xfa\x18\x25\xea\xda\x92\x91\x95\x26\x69\xc9\xbd\x97\x91\xd3\x58\x0b\xed\xfe\x09\xf8\x39\x33\x73\xe6\x4c\xc4\x83\xf7\x11\x7e\x15\x00\x55\xe5\x7c\xd0\x35\x87\xb2\xc1\xb6\xeb\xea\xf7\x27\x63\x1c\xca\x56\x76\x4d\xdf\xec\xac\xe1\x50\x76\x6d\xa7\x7a\xdc\x59\x4b\x4c\xf5\xc3\x20\x76\xd6\x71\x28\xf5\xa8\x11\x0f\x3b\xeb\x39\x94\xd8\xe3\xc1\x64\x1e\x03\x31\x85\xc6\x74\x3b\x7b\xe3\x50\x8e\x46\x2a\x29\x77\x36\x12\x1b\x55\xad\xb3\xde\x03\x31\x26\x98\x62\xd9\xce\x14\xa4\xc7\x91\x09\x95\x41\x4a\x22\xcd\xc0\xf2\x0d\x19\x45\xd1\xf5\xf8\xf6\x96\x47\xa6\x2c\x28\x95\x1c\x33\x6f\x46\x61\x44\x2b\x71\xcc\x67\x52\x1a\xd9\x8d\x28\xf4\x7b\x71\x2f\x8a\x6f\xe9\x98\xca\x4f\x3e\x70\x28\x8d\x31\x54\x6a\xbc\x8b\x95\x11\xb3\x9d\x6e\x1c\x66\xef\xfc\xba\x08\x85\x4f\xe5\x82\xf6\x78\x8a\x1c\x9c\x0f\xb3\x98\x08\xcf\x22\x1c\xad\xe3\x90\x56\x5a\x84\xd6\xd6\x1d\xd3\xd7\xbd\x28\x4e\x71\x9e\x92\x89\x14\xea\xe3\x18\xfc\xd9\xe9\xea\xe1\xf7\x43\x84\xff\x1f\xff\xe3\x0b\x75\x9e\x1e\x83\x59\x5d\xff\x97\x7a\xa5\xd7\xb7\xd4\xfb\x59\xf9\xc2\x10\x40\xdb\x75\x99\xc4\x8d\x83\x99\xf0\x4a\x40\x4c\xf6\xe8\x2a\xe5\x5d\x44\x17\x39\x28\x74\x11\x03\x09\xdf\xcf\x6b\xb4\xe6\xf6\x37\x69\xeb\xb1\x11\xe7\x75\x1b\x54\xad\x51\x84\x98\x4e\x40\x5f\x97\x20\x16\x0e\xf4\x4b\x28\xe2\x35\x56\xa9\x65\x9f\x71\x2f\x0a\xe3\xc3\x9c\xb6\xff\xb3\x2d\xab\x97\x6b\x52\xac\x5b\xce\xdb\x33\x96\x3e\x68\x0c\x74\x4b\x87\x9f\xb2\x6c\xb5\x00\x17\xab\xe3\x89\x43\xbf\xe5\xfe\xea\x86\xec\x25\xcd\x7e\xa5\xad\x2b\xd6\xa4\xf1\xff\x1a\xf7\xbc\xf0\x56\x3d\xa4\xe2\xdc\xea\x29\xb5\xb9\xd4\xf4\xb9\xd4\xe4\x12\x1b\x5e\x87\x4d\x2c\xd5\xe4\xd5\x47\xd2\xd2\x73\x59\xed\x4f\x24\xc7\x80\x33\xe9\xbf\x03\x00\x00\xff\xff\x48\x0e\x7d\x47\xc3\x03\x00\x00")

func assetsStylesCssBytes() ([]byte, error) {
	return bindataRead(
		_assetsStylesCss,
		"assets/styles.css",
	)
}

func assetsStylesCss() (*asset, error) {
	bytes, err := assetsStylesCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/styles.css", size: 963, mode: os.FileMode(420), modTime: time.Unix(1575982167, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xbf\x4e\xf4\x30\x10\xc4\x6b\xe7\x29\xfc\x6d\xff\x9d\x15\x6a\x3b\x0d\x87\x28\xa1\xa0\xa1\x34\xeb\x3d\x6c\xe2\xd8\xc1\xeb\x20\xee\xed\x91\xcf\x3a\xf1\xa7\x38\x09\x8a\x55\x8a\x99\xdf\xc4\x9a\xd1\xff\xf6\x77\xd7\x0f\x8f\xf7\x37\xd2\xd7\x25\x4e\x83\xee\x1f\xa1\x3d\x59\x37\x0d\x42\xe8\x18\xd2\x2c\x0b\x45\x03\x5c\x8f\x91\xd8\x13\x55\x90\xbe\xd0\xc1\x80\xb2\xcc\x54\x59\x75\x65\x87\xcc\xd0\x58\xd5\xe1\x41\xe8\xa7\xec\x8e\xa7\x14\x3f\xca\xe0\x0c\x60\xcc\x38\x83\xc4\x68\x99\x0d\x84\x4a\xcb\xff\xf1\x0a\x26\xad\xfc\xf8\xd5\xe6\x6c\xa5\x0b\xae\x43\x2e\xcb\x4f\x55\x26\xbb\x90\x81\x77\x90\x16\x6b\xc8\xc9\x80\x52\x6e\xc3\xb9\xdd\x73\xde\x61\x5e\x54\x7b\x9b\x10\x3a\xa4\x75\xab\xb2\x1e\x57\x32\xe0\x83\x73\x94\x40\xbe\xd9\xb8\x91\x81\xf1\x1c\x33\xfb\xf6\xbb\x93\xf3\x17\x50\xfa\x0b\x64\xf1\x02\xc5\x64\x0b\x7a\x90\x6b\xb4\x48\x3e\x47\x47\xc5\xc0\x7e\xc3\xb9\xdd\x6d\x3e\x87\xbc\x7e\x8b\xd0\xaa\x15\x74\x1a\xa2\xf7\x3f\x08\xcd\x58\xc2\x5a\x25\x17\xfc\x5c\xad\xb5\xbc\x7b\xe1\xc6\x76\x79\x1a\xb4\xea\xfb\x7f\x04\x00\x00\xff\xff\xbe\x74\xd0\x85\x17\x02\x00\x00")

func htmlIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlIndexHtml,
		"html/index.html",
	)
}

func htmlIndexHtml() (*asset, error) {
	bytes, err := htmlIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/index.html", size: 535, mode: os.FileMode(420), modTime: time.Unix(1575982167, 0)}
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
	"assets/date.js":    assetsDateJs,
	"assets/styles.css": assetsStylesCss,
	"html/index.html":   htmlIndexHtml,
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
	"assets": &bintree{nil, map[string]*bintree{
		"date.js":    &bintree{assetsDateJs, map[string]*bintree{}},
		"styles.css": &bintree{assetsStylesCss, map[string]*bintree{}},
	}},
	"html": &bintree{nil, map[string]*bintree{
		"index.html": &bintree{htmlIndexHtml, map[string]*bintree{}},
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

func assetFS() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}
