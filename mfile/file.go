package mfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/m-startgo/go-utils/mpath"
)

// ReadFile 读取文件内容为字节切片
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile 将字节写入文件（覆盖），会尝试创建父目录
func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if dir != "" {
		if _, err := mpath.EnsureDir(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, perm)
}

// Base 返回文件名（带扩展名）
func Base(p string) string {
	return filepath.Base(p)
}

// Ext 返回文件扩展名（包含点）
func Ext(p string) string {
	return filepath.Ext(p)
}

// HasExt 判断文件是否具有给定扩展（不区分大小写）
func HasExt(p, ext string) bool {
	return strings.EqualFold(filepath.Ext(p), ext)
}

// ListFiles 列出目录下（不递归）的文件名（不包含目录）
func ListFiles(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(fis))
	for _, fi := range fis {
		if !fi.IsDir() {
			out = append(out, fi.Name())
		}
	}
	return out, nil
}
