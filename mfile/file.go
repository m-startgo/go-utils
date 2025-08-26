package mfile

import (
	"os"
	"path/filepath"

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
