// Code generated by go-bindata.
// sources:
// config.json
// DO NOT EDIT!

package config

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _configJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd4\x90\xc1\x6e\x84\x20\x14\x45\xf7\x7c\x05\x61\xed\x0c\xfb\xf9\x95\xa6\x0b\x0b\x2f\x4a\xd4\xf7\x08\xa0\xb1\x69\xfc\xf7\x22\xd4\x16\x4d\x9a\xd8\x4d\xe3\xec\xcc\xbd\xcf\xcb\xc9\xf9\x60\x9c\x0b\x40\x6d\xc9\x60\xf0\xe2\xc1\xd7\x20\x46\x61\x44\x84\xde\x83\x9b\xc0\xc5\xf4\x25\xa5\xfc\xab\x4d\x17\x80\x93\x71\x84\x03\x60\x28\x0e\x52\x65\x1d\xe9\x51\x05\x43\x28\xaa\x32\x1f\x6a\xac\x1b\xd0\xe2\x3b\x7b\xfd\xa9\xc5\xe8\xfa\xb8\x23\xda\x10\xec\x43\xca\xfc\x7e\x1c\x9a\xdf\xef\x1d\x69\x83\xcd\x5d\xd1\x20\x3b\x13\x60\xfb\x7d\xa9\xfe\x04\xa5\x61\x82\x9e\x6c\xaa\xaa\x43\xb1\x81\xed\x72\x5f\xa3\x7e\xa3\xf9\x0c\x6d\x9c\x38\x03\xcc\x8a\x15\x61\xec\xbf\x7b\xf5\x11\xd5\x96\x74\x37\x19\x29\x2e\xa6\xd3\x67\x9f\xb7\xdf\x41\x0f\x1a\x55\x0b\xaa\xbb\x86\xcb\x8c\xf2\x2c\x42\xf7\xb4\xd9\x2a\x5b\xbf\x16\xf6\x19\x00\x00\xff\xff\x7a\x20\xdf\x50\x17\x04\x00\x00")

func configJsonBytes() ([]byte, error) {
	return bindataRead(
		_configJson,
		"config.json",
	)
}

func configJson() (*asset, error) {
	bytes, err := configJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.json", size: 1047, mode: os.FileMode(420), modTime: time.Unix(1446555960, 0)}
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
	"config.json": configJson,
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
	"config.json": &bintree{configJson, map[string]*bintree{}},
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

