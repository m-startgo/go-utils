package mpath

import (
	"os"
)

// Home 返回用户的主目录。
// 它会尝试跨平台常用的环境变量。
// 若均未设置，返回空字符串。
func Home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return dir
}

// Pwd 返回当前工作目录及可能发生的错误。
// 建议调用方处理返回的 error，而不是吞掉它。
func Pwd() (string, error) {
	return os.Getwd()
}

// 判断目录或文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

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
