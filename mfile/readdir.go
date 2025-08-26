package mfile

import (
	"errors"
	"os"
	"path/filepath"
)

// FileNode 表示目录中每一项的返回信息
type FileNode struct {
	Name       string `json:"name"`     // 文件或目录名
	RelPath    string `json:"rel_path"` // 相对路径，相对于传入的根目录
	AbsPath    string `json:"abs_path"` // 绝对路径
	IsFile     bool   `json:"is_file"`  // 是否为文件
	DirName    string `json:"dir_name"` // 所在目录名
	DirRelPath string `json:"dir_rel"`  // 所在目录相对路径
	DirAbsPath string `json:"dir_abs"`  // 所在目录绝对路径
	IsDir      bool   `json:"is_dir"`   // 是否为目录
}

// ReadDir 列出目录下的文件和目录。level 表示递归深度：
//
//	 0 => 只列出当前目录（不递归）
//	>0 => 递归到指定层级（相对于根目录）
//	-1 => 递归所有层级
//
// 返回每项的 FileNode 列表
func ReadDir(root string, level int) ([]FileNode, error) {
	if root == "" {
		return nil, errors.New("root path empty")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("root path is not a directory")
	}
	var res []FileNode

	var walk func(current string, currentLevel int) error
	walk = func(current string, currentLevel int) error {
		entries, err := os.ReadDir(current)
		if err != nil {
			return err
		}
		for _, e := range entries {
			name := e.Name()
			absPath := filepath.Join(current, name)
			relPath, _ := filepath.Rel(absRoot, absPath)
			parentDir := filepath.Dir(absPath)
			parentRelDir, _ := filepath.Rel(absRoot, parentDir)

			node := FileNode{
				Name:       name,
				RelPath:    relPath,
				AbsPath:    absPath,
				IsFile:     !e.IsDir(),
				DirName:    filepath.Base(parentDir),
				DirRelPath: parentRelDir,
				DirAbsPath: parentDir,
				IsDir:      e.IsDir(),
			}
			res = append(res, node)

			// 如果需要递归且当前是目录
			if e.IsDir() {
				// decide whether to go deeper
				if level == -1 || currentLevel < level {
					if err := walk(absPath, currentLevel+1); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	if err := walk(absRoot, 0); err != nil {
		return nil, err
	}
	return res, nil
}
