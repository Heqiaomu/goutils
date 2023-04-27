package compress

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	util "github.com/Heqiaomu/goutil"
)

// IsZip return true if file is zip
func IsZip(zipPath string) bool {
	if exist, _ := util.PathExists(zipPath); !exist {
		fmt.Println("zip file not found: ", zipPath)
		return false
	}
	f, err := os.Open(zipPath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

// Unzip 解压zip文件
func Unzip(archive, target string) (err error) {
	if exist, _ := util.PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	if !IsZip(archive) {
		return fmt.Errorf("target archive <%s> is not a zip file", archive)
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if strings.Contains(path, "__MACOSX") {
			continue
		}
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return
}

// UnzipWithFileName 解压zip文件，将加压后的文件解压到指定文件名的文件中
// 例如 archive = "./test_data/drivers/hpcv1_0_5.zip"
//     target = "./test_data/drivers/"
//     fileName = "31263478912"
// 那么最终将会在 ./test_data/drivers/31263478912 文件下生成所有原本解压在 hpcv1_0_5 文件夹下的所有内容（不包含hpcv1_0_5这个文件夹）
func UnzipWithFileName(archive, target string, fileName string) (err error) {
	if exist, _ := util.PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	if !IsZip(archive) {
		return fmt.Errorf("target archive <%s> is not a zip file", archive)
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(target, fileName), 0755); err != nil {
		return err
	}

	// 有三种情况
	// 1. 没有根文件夹 (只有 . || 没有 . ) && (不包含 / || \)
	// 2. 有根文件夹，但是不是单独的File,
	// 3. 有根文件夹，也是单独的File, 有 / 或者 \ 结尾
	var first bool
	//var second = true
	var third bool
	var prefix string
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "/") || strings.HasSuffix(file.Name, "\\") {
			third = true
			prefix = file.Name
			break
		}
		if !(strings.Contains(file.Name, "/") || strings.Contains(file.Name, "\\")) &&
			(strings.Contains(file.Name, ".") || len(file.Name) != 0) {
			first = true
			break
		}
	}

	if !first && !third {
		var index int
		if strings.HasPrefix(reader.File[0].Name, "/") ||
			strings.HasPrefix(reader.File[0].Name, "\\") {
			index = 1
		}
		split := strings.Split(reader.File[0].Name, "/")
		if len(split) == 1 {
			split = strings.Split(reader.File[0].Name, "\\")
			if index != 0 {
				prefix = "\\" + split[index] + "\\"
			} else {
				prefix = split[index] + "\\"
			}
		} else {
			if index != 0 {
				prefix = "/" + split[index] + "/"
			} else {
				prefix = split[index] + "/"
			}
		}
		if len(prefix) == 0 {
			return fmt.Errorf("fail to check prefix, because length is 0. please check file [%s] is the correct zip format", archive)
		}
	}

	// rootName = hpcv1.0.5/
	//rootName := ""
	for _, file := range reader.File {
		// 第一次获取的是文件夹的名字
		//if i == 0 {
		//	rootName = file.Name
		//}

		oldFileName := file.Name
		// 跳过os系统压缩的文件
		if strings.Contains(oldFileName, "__MACOSX") {
			continue
		}

		//if !strings.HasPrefix(oldFileName, rootName) {
		//	return fmt.Errorf("fail to check file name, because not contains root file name. please compressed file must in a root file while file=%s not in root file=%s", oldFileName, rootName)
		//}

		// 路径中的根文件夹的名字替换成指定的文件夹名
		var newFileName string
		if first {
			newFileName = filepath.Join(fileName, oldFileName)
		} else {
			newFileName = filepath.Join(fileName, strings.TrimPrefix(oldFileName, prefix))
		}
		newFileName = strings.ReplaceAll(newFileName, "\\", "/")
		//newFileName := filepath.Join(fileName, oldFileName[len(rootName):])
		path := filepath.Join(target, newFileName)

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return
}

// Zipdir  压缩目录
// src     需要压缩的目录
// dest    压缩文件存放路径
func Zipdir(src string, dest string) (err error) {
	if exist, _ := util.PathExists(src); !exist {
		return fmt.Errorf("src path not found: %s", src)
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	return Zipfiles([]*os.File{f}, dest)
}

// Zipfiles 压缩文件
// files    文件数组
// dest     压缩文件存放路径
func Zipfiles(files []*os.File, dest string) (err error) {
	// subfix .zip
	if !strings.HasSuffix(dest, ".zip") {
		dest = strings.Join([]string{dest, ".zip"}, "")
	}

	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err = zipfile(file, "", w)
		if err != nil {
			return err
		}
	}

	return
}

func zipfile(file *os.File, prefix string, zw *zip.Writer) (err error) {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = zipfile(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}

	return
}
