package mpath

import (
	"os"
)

// 判断目录或文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(nil) {
		return false
	}
	return false
}
