package mpath

import "os"

// IsDir 判断给定路径是否为目录。
// 对于不存在的路径或发生错误时返回 false。
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 判断给定路径是否存在且为常规文件。
// 对于不存在的路径、目录或发生错误时返回 false。
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
