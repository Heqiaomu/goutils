package util

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	// FilePermission file mode
	FilePermission = 0644
)

// GetDirList get all dir
func GetDirList(dirpath string) ([]string, error) {
	var dirList []string
	files, err := ioutil.ReadDir(dirpath)
	for _, fi := range files {
		if fi.IsDir() {
			dirList = append(dirList, filepath.Join(dirpath, fi.Name()))
		}
	}
	return dirList, err
}

// RemoveDirContents remove dir contents
func RemoveDirContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// FileExist judge file exist
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// Gen generate file with template and data
func Gen(data interface{}, tpl, filepath string) (err error) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return
	}

	err = WriteFile(buf.Bytes(), filepath)
	if err != nil {
		return
	}

	return nil
}

// GenR gen file with template and return
func GenR(data interface{}, tplT []byte) ([]byte, error) {
	t := template.New("temp")
	t, err := t.Parse(string(tplT))
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

//AdvanceGenR add indent support
func AdvanceGenR(data interface{}, tplT []byte) ([]byte, error) {
	t := template.New("temp").Funcs(sprig.TxtFuncMap())
	t, err := t.Parse(string(tplT))
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// WriteFile will backup file for rollback when write fail
func WriteFile(data []byte, f string) error {
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE, FilePermission)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// Ext return file extension
// e.g. Ext("abc.json") => "json"
func Ext(filename string, checkList ...string) (ext string, err error) {
	ext = filepath.Ext(filename)
	if len(ext) <= 1 {
		return ext, fmt.Errorf("filename: %s requires valid extension", filename)
	}

	ext = ext[1:]
	if len(checkList) > 0 && !StringInSlice(ext, checkList) {
		return ext, fmt.Errorf("Unsupported Config Type %s", ext)
	}

	return ext, nil
}

// CreatedDirIfNotExist 创建文件夹如果不存在
func CreatedDirIfNotExist(dir string) error {
	exist, err := PathExists(dir)
	if err != nil {
		return err
	}
	if !exist {
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
