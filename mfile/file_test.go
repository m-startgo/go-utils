package mfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteRead(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "a.txt")
	content := "hello mfile"
	if err := Write(fp, content); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	b, err := Read(fp)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if string(b) != content {
		t.Fatalf("content mismatch: got %q want %q", string(b), content)
	}
}

func TestAppend(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "b.txt")
	if err := Write(fp, "one"); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err := Append(fp, "-two"); err != nil {
		t.Fatalf("Append error: %v", err)
	}
	b, err := Read(fp)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if string(b) != "one-two" {
		t.Fatalf("append mismatch: got %q", string(b))
	}
}

func TestDetectMimeAndExt(t *testing.T) {
	// PNG header
	png := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	mt := DetectMime(png)
	if mt == "" {
		t.Fatalf("DetectMime failed for png")
	}
	ext := ExtByContent(png)
	if ext == "" {
		t.Fatalf("ExtByContent failed for png")
	}
}

func TestReadDir(t *testing.T) {
	dir := t.TempDir()
	// create structure:
	// dir/
	//   f1.txt
	//   .hidden
	//   sub/
	//     f2.txt
	os.WriteFile(filepath.Join(dir, "f1.txt"), []byte("1"), 0o644)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("h"), 0o644)
	os.Mkdir(filepath.Join(dir, "sub"), 0o755)
	os.WriteFile(filepath.Join(dir, "sub", "f2.txt"), []byte("2"), 0o644)

	// level 0 should list only items directly under dir
	nodes, err := ReadDir(dir, 0)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 entries in level0, got %d", len(nodes))
	}

	// level -1 should list all including sub/f2.txt
	nodesAll, err := ReadDir(dir, -1)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	// expecting 4 entries: f1, .hidden, sub, sub/f2
	if len(nodesAll) != 4 {
		t.Fatalf("expected 4 entries in level-1, got %d", len(nodesAll))
	}
}

func TestEmptyPathErrors(t *testing.T) {
	// Write / WriteByte
	if err := Write("", "x"); err == nil {
		t.Fatalf("expected error for empty path in Write")
	}
	if err := WriteByte("", []byte("x")); err == nil {
		t.Fatalf("expected error for empty path in WriteByte")
	}

	// Append / AppendByte
	if err := Append("", "x"); err == nil {
		t.Fatalf("expected error for empty path in Append")
	}
	if err := AppendByte("", []byte("x")); err == nil {
		t.Fatalf("expected error for empty path in AppendByte")
	}

	// Read
	if _, err := Read(""); err == nil {
		t.Fatalf("expected error for empty path in Read")
	}

	// ReadDir
	if _, err := ReadDir("", 0); err == nil {
		t.Fatalf("expected error for empty path in ReadDir")
	}
}

func TestReadDirOrderingAndHidden(t *testing.T) {
	dir := t.TempDir()
	// create files with varying names
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(dir, ".z"), []byte("z"), 0o644)

	nodes, err := ReadDir(dir, 0)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(nodes))
	}
	// Expect alphabetical order by name (including dotfiles where '.' < 'a')
	if nodes[0].Name != ".z" || nodes[1].Name != "a.txt" || nodes[2].Name != "b.txt" {
		t.Fatalf("unexpected order: %#v", nodes)
	}
}

func TestWriteAppendErrors(t *testing.T) {
	// 在所有平台上可靠失败的情况：尝试向已存在的目录路径写入内容
	dir := t.TempDir()
	targetDir := filepath.Join(dir, "targetdir")
	if err := os.Mkdir(targetDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	// 尝试写入一个已经存在且为目录的路径，应返回错误
	if err := Write(targetDir, "x"); err == nil {
		t.Fatalf("expected error when writing to a path that is a directory")
	}
	if err := Append(targetDir, "x"); err == nil {
		t.Fatalf("expected error when appending to a path that is a directory")
	}
}

func TestDetectMimeEmpty(t *testing.T) {
	if DetectMime(nil) != "" {
		t.Fatalf("expected empty mime for nil content")
	}
	if ExtByContent(nil) != "" {
		t.Fatalf("expected empty ext for nil content")
	}
}
