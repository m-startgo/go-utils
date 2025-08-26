package mpath

import (
	"os"
	"path/filepath"
)

// Join 将路径片段安全拼接为一个路径
// Join("a", "b", "c.txt") -> "a/b/c.txt"
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Abs 返回给定路径的绝对路径
// Abs("./a") -> ("/full/path/to/a", nil)
func Abs(p string) (string, error) {
	return filepath.Abs(p)
}

// EnsureDir 确保目录存在，若目录不存在则尝试创建。
// 参数:
//
//	dir  - 目标目录路径；若为空字符串，按仓库约定返回 (false, nil)
//	perm - 创建目录时使用的权限掩码
//
// 返回:
//
//	created - 如果本次调用创建了目录返回 true；如果目录已存在返回 false
//	err     - 发生错误时返回非 nil（例如权限问题或创建失败）
//
// 说明:
//   - 对于并发创建场景，如果 MkdirAll 失败但随后检测到目录已存在，会视为创建成功并返回 true。
//
// 用例:
//
//	EnsureDir("/tmp/a/b", 0755) -> (true, nil) // 如果目录之前不存在
func EnsureDir(dir string, perm os.FileMode) (bool, error) {
	// 空路径：遵循仓库约定，返回零值而非错误
	if dir == "" {
		return false, nil
	}

	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return false, nil // 已存在，未新建
	}
	if err == nil && !info.IsDir() {
		return false, &os.PathError{Op: "mkdir", Path: dir, Err: os.ErrExist}
	}

	// 如果不存在，尝试创建；若创建失败，则再次 stat 以处理并发创建的情况
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, perm); err != nil {
			if info2, statErr := os.Stat(dir); statErr == nil && info2.IsDir() {
				return true, nil
			}
			return false, err
		}
		return true, nil
	}
	// 其它错误（例如权限问题）直接返回
	return false, err
}

// Rel 返回从 base 到 p 的相对路径（filepath.Rel 的包装）
// 参数:
//
//	p    - 目标路径
//	base - 基准路径
//
// 返回: 相对路径字符串，以及可能的错误（若无法计算相对路径）
// 用例: Rel("/a/b/c", "/a") -> ("b/c", nil)
func Rel(p, base string) (string, error) {
	return filepath.Rel(base, p)
}

// Clean 规范化路径，移除多余的分隔符和相对段（wrapper for filepath.Clean）
// Clean("/a//b/./c") -> "/a/b/c"
func Clean(p string) string {
	return filepath.Clean(p)
}

// IsAbs 判断路径是否为绝对路径（wrapper for filepath.IsAbs）
// IsAbs("/a/b") -> true
func IsAbs(p string) bool {
	return filepath.IsAbs(p)
}
