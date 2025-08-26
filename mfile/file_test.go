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
