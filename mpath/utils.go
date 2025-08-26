package mpath

import (
	"os"
	"path/filepath"
)

// Join 将路径片段安全拼接为一个路径
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Abs 返回给定路径的绝对路径
func Abs(p string) (string, error) {
	return filepath.Abs(p)
}

// EnsureDir 确保目录存在（若不存在则创建），返回是否新创建或错误
func EnsureDir(dir string, perm os.FileMode) (bool, error) {
	if dir == "" {
		return false, nil
	}
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return false, nil
	}
	if err == nil && !info.IsDir() {
		return false, &os.PathError{Op: "mkdir", Path: dir, Err: os.ErrExist}
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, perm); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}
