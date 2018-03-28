package template

import (
	"bytes"
	"compress/gzip"
	"fmt"
	tmplhtml "html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _templateDefaultTmpl = []byte("\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\x8c\x91\xc1J\xf3P\x10\x85\xf7y\x8aC\xe9\xf2O\xdbu\xc9\x1f\xe8J\x04w\xe2JDn\xec\xb4\x06o'5\xb9)\xc88K\xdf\xc0\xa5\xe0\xca\x17\xf3\t|\x04\xb9\xb9\x89M\x83\xd4\xee\x86ə\x93s\xbe+\x82\xb1%^\xbb{\xcc\xff\xc3\x12crV\x16\xf5\xf6\xc2dd+Ī\x10A\xbe\x02ӏp\x06\xd5(\xc9\xd2\xebF\x89 \xbdI\xa6Y\x1a\x89\xa04\xbc&\x8c\x1f\xe8\xe9\xdfxg\xac\xb7\x1dZF\x00\x90\xd0&\x15\xf1\x12\xd5dJ\x1b\u007f\x1c\x83x\tUD\"\xf9\n\xf4\x88ɥ3\xae\xae0Z\xe5e\xce\xeb\x91j\x146s$Y*\xd2}\u007f\x86+\xae\xb6[*\xfd\xf5\xd7\xfb\xebG\x1b\x87x\xa9\xda8\x0f\xfcJ\xaa\n\xbb\xa3\xe5ߎ\xaa\xf8|{9\xf0K\xb2\xf4\x9c+g\xf8\x8e\xe6\x83\xde]入҅\xb6{Ʒ60hQ{\xf1\xe4\b\xe9N>kl<\xf1\xe3\xack\xf2\xce\x03\xd7\x00\xdb7\xf3:\xd5&0\x02\xf9\x9aB\xbc=\xf5f\x8c۹\x8ba\x98\vg\\^\xf0a\xf4Eo\xffk\xfe\xfe\xe1\xbeD\xef\xec\x84&\x83\x9fD\xa1\rN\xa8\x03\xf4\x1b\x85\xc7k\xa7\xb8\xdb~\x03\x00\x00\xff\xff\x01\x00\x00\xff\xffu^\x86\x11\x00\x03\x00\x00")

func bindataWriter(tmp []byte) ([]byte, error) {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(tmp)
	if err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func TempateRead(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	tapmlate, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bindataWriter(tapmlate)
}

func BindataRead(data []byte, name string) ([]byte, error) {
	return bindataRead(data, name)
}

func getDefaultTemplate(name string, mapsInstance map[string]string) (*tmplhtml.Template, error) {
	tempateData, err := bindataRead(_templateDefaultTmpl, name)
	if err != nil {
		return nil, err
	}
	return tmplhtml.New(name).Option("missingkey=zero").Funcs(initFuncMap(mapsInstance)).Parse(string(tempateData))
}

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

func templateDefaultTmplBytes() ([]byte, error) {
	return bindataRead(
		_templateDefaultTmpl,
		"template/default.tmpl",
	)
}

func templateDefaultTmpl() (*asset, error) {
	bytes, err := templateDefaultTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/default.tmpl", size: 768, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
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
	"template/default.tmpl": templateDefaultTmpl,
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
	"template": &bintree{nil, map[string]*bintree{
		"default.tmpl": &bintree{templateDefaultTmpl, map[string]*bintree{}},
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
