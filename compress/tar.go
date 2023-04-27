package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	util "github.com/Heqiaomu/goutil"
)

// Tardir 压缩目录为tar.gz
// @src : if src == "../app", then Untar the tar.gz, the first fold name is "app"
// @dest: must contains the compress file's name, such as "./test.tar.gz"
func Tardir(src, dest string) (err error) {
	if exist, _ := util.PathExists(src); !exist {
		return fmt.Errorf("src path not found: %s", src)
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	return Tarfile([]*os.File{f}, dest)
}

// Tarfile 压缩gzip压缩成tar.gz
func Tarfile(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := tarfile(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func tarfile(file *os.File, prefix string, tw *tar.Writer) error {
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
			err = tarfile(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Untar 解压 tar.gz
// @archive : the compress file's path , such as - "./test.tar.gz"
// @dest    : the file Untar dest, such as "./ttt", and you will get fold named "test" under the ttt
// @force   : {false} -> will check dest path exist, {true} -> not check,and will clean all file under dest
func Untar(archive, dest string, force bool) error {
	if exist, _ := util.PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	if !force {
		if e, _ := util.PathExists(dest); e {
			return fmt.Errorf("dest folder alreay exits: %s", dest)
		}
	}

	// remove dest
	if e, _ := util.PathExists(dest); e {
		err := os.RemoveAll(dest)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := path.Join(dest, hdr.Name) // dest + hdr.Name

		if hdr.FileInfo().IsDir() {
			// create path before create file in <create> func, continue here
			continue
		}

		file, err := create(filename)
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			_ = file.Close()
			return err
		}
		_ = file.Close()
		_ = os.Chmod(filename, hdr.FileInfo().Mode())
	}
	return nil
}

// Untargz 解压 tar.gz
// @archive : the compress file's path , such as - "./test.tar.gz"
// @dest    : the file Untar dest, such as "./ttt", and you will get fold named "test" under the ttt
func Untargz(archive, dest string) error {
	if exist, _ := util.PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	// 如果不存在，则新建dest
	if e, _ := util.PathExists(dest); !e {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := path.Join(dest, hdr.Name) // dest + hdr.Name

		if hdr.FileInfo().IsDir() {
			// create path before create file in <create> func, continue here
			continue
		}

		file, err := create(filename)
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			_ = file.Close()
			return err
		}
		_ = os.Chmod(filename, hdr.FileInfo().Mode())
		_ = file.Close()
	}
	return nil
}

// Untargz decompression tar.gz
// @archive : the compress file's path , such as - "./test.tar.gz"
// @dest    : the file Untar dest, such as "./ttt", and you will get fold named "test" under the ttt
// @fileName : the new name for the extracted file. If your compressed file is xxx.tar.gz which contains /xxx/ss.text, and your dest=./driver and fileName=aaa, then you will get /driver/aaa/ss.text
func UntargzWithName(archive, dest string, fileName string) error {
	if exist, _ := util.PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	// 如果不存在，则新建dest
	if e, _ := util.PathExists(dest); !e {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	hdr, err := tr.Next()
	if err != nil {
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	if !hdr.FileInfo().IsDir() {
		// create path before create file in <create> func, continue here
		return fmt.Errorf("compressed file must contains a dir in the root path")
	}

	rootName := hdr.Name

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		hdrName := hdr.Name
		if !strings.HasPrefix(hdrName, rootName) {
			return fmt.Errorf("compressed file must in a root file while file=%s not in root file=%s", hdrName, rootName)
		}
		hdrName = path.Join(fileName, hdrName[len(rootName):])

		filename := path.Join(dest, hdrName) // dest + hdr.Name

		if hdr.FileInfo().IsDir() {
			// create path before create file in <create> func, continue here
			continue
		}

		file, err := create(filename)
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			_ = file.Close()
			return err
		}
		_ = os.Chmod(filename, hdr.FileInfo().Mode())
		_ = file.Close()
	}

	return nil

}

func create(name string) (*os.File, error) {
	dir, _ := filepath.Split(name)
	// create dir before create file
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
