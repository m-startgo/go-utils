package mfile

import (
	"errors"
	"os"
	"path/filepath"
)

// Write 将字符串内容写入文件，若目录不存在则创建，若文件存在则覆盖
func Write(filePath string, content string) error {
	return WriteByte(filePath, []byte(content))
}

// WriteByte 将字节内容写入文件，若目录不存在则创建，若文件存在则覆盖
func WriteByte(filePath string, content []byte) error {
	if filePath == "" {
		return errors.New("file path empty")
	}
	// 标准化路径，避免出现类似 "a/../b" 导致的异常
	filePath = filepath.Clean(filePath)
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, content, 0o644)
}

// Append 将字符串追加到文件末尾，若文件不存在则创建
func Append(filePath string, content string) error {
	return AppendByte(filePath, []byte(content))
}

// AppendByte 将字节追加到文件末尾，若文件不存在则创建
func AppendByte(filePath string, content []byte) (err error) {
	if filePath == "" {
		return errors.New("file path empty")
	}
	filePath = filepath.Clean(filePath)
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	// 确保 Close 的错误不会被忽略
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	n, werr := f.Write(content)
	if werr != nil {
		return werr
	}
	if n != len(content) {
		return errors.New("write incomplete")
	}
	return nil
}

// Read 返回文件内容的字节以及错误
func Read(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, errors.New("file path empty")
	}
	return os.ReadFile(filePath)
}
