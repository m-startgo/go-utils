package mpath

import (
	"os"
	"runtime"
)

// Home 返回用户的主目录。
// 它会尝试跨平台常用的环境变量。
// 若均未设置，返回空字符串。
func Home() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Pwd 返回当前工作目录及可能发生的错误。
// 建议调用方处理返回的 error，而不是吞掉它。
func Pwd() (string, error) {
	return os.Getwd()
}
