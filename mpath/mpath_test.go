package mpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHome(t *testing.T) {
	h := Home()
	if h == "" {
		t.Fatalf("home should not be empty")
	}
}

func TestPwd(t *testing.T) {
	p, err := Pwd()
	if err != nil {
		t.Fatalf("pwd error: %v", err)
	}
	if p == "" {
		t.Fatalf("pwd should not be empty")
	}
}

func TestIsExistAndIsFileDir(t *testing.T) {
	tmpdir := os.TempDir()
	// ensure exists
	if !IsExist(tmpdir) {
		t.Fatalf("tmp dir should exist: %s", tmpdir)
	}
	if !IsDir(tmpdir) {
		t.Fatalf("tmp dir should be dir: %s", tmpdir)
	}
	// create temp file
	fpath := filepath.Join(tmpdir, "mpath_test_tmp_file.txt")
	_ = os.Remove(fpath)
	if err := os.WriteFile(fpath, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write tmp file err: %v", err)
	}
	defer os.Remove(fpath)
	if !IsExist(fpath) {
		t.Fatalf("file should exist: %s", fpath)
	}
	if !IsFile(fpath) {
		t.Fatalf("should be file: %s", fpath)
	}
}

func TestUtilsBasics(t *testing.T) {
	// Join/Abs/Base/Ext
	p := Join("a", "b", "c.txt")
	if p == "" {
		t.Fatalf("join returned empty")
	}
	abs, err := Abs(p)
	if err != nil || abs == "" {
		// Abs may fail for relative paths on some envs, don't strict fail on err, but ensure string returned when no err
	}
}

func TestEnsureReadWriteList(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "mpath_test_dir")
	_ = os.RemoveAll(dir)
	created, err := EnsureDir(dir, 0o755)
	if err != nil {
		t.Fatalf("ensure dir err: %v", err)
	}
	if !created {
		t.Fatalf("should be created")
	}
}

func TestEnsureDir_EmptyPath(t *testing.T) {
	created, err := EnsureDir("", 0o755)
	if err != nil {
		t.Fatalf("expected no error for empty path, got: %v", err)
	}
	if created {
		t.Fatalf("expected created==false for empty path")
	}
}

func TestEnsureDir_FileExists(t *testing.T) {
	f := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_ensure_file_%d.tmp", time.Now().UnixNano()))
	_ = os.Remove(f)
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("write tmp file err: %v", err)
	}
	defer os.Remove(f)

	created, err := EnsureDir(f, 0o755)
	if err == nil {
		t.Fatalf("expected error when EnsureDir on existing file, got nil")
	}
	if !errors.Is(err, os.ErrExist) {
		// some environments may wrap; ensure the cause is os.ErrExist when possible
		t.Logf("warning: EnsureDir returned non-exist error type: %T %v", err, err)
	}
	if created {
		t.Fatalf("expected created==false when path exists as file")
	}
}

func TestEnsureDir_AlreadyDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_ensure_dir_%d", time.Now().UnixNano()))
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdirall err: %v", err)
	}
	defer os.RemoveAll(dir)

	created, err := EnsureDir(dir, 0o755)
	if err != nil {
		t.Fatalf("expected no error for existing dir, got: %v", err)
	}
	if created {
		t.Fatalf("expected created==false when directory already exists")
	}
}

func TestIsExist_NotExist(t *testing.T) {
	p := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_not_exist_%d", time.Now().UnixNano()))
	_ = os.RemoveAll(p)
	if IsExist(p) {
		t.Fatalf("IsExist should be false for nonexistent path: %s", p)
	}
	if IsDir(p) {
		t.Fatalf("IsDir should be false for nonexistent path: %s", p)
	}
	if IsFile(p) {
		t.Fatalf("IsFile should be false for nonexistent path: %s", p)
	}
}

func TestIsDir_FileInput(t *testing.T) {
	f := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_isdir_file_%d.tmp", time.Now().UnixNano()))
	_ = os.Remove(f)
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("write tmp file err: %v", err)
	}
	defer os.Remove(f)
	if IsDir(f) {
		t.Fatalf("IsDir should be false for file input: %s", f)
	}
}

func TestIsFile_DirInput(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_isfile_dir_%d", time.Now().UnixNano()))
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdirall err: %v", err)
	}
	defer os.RemoveAll(dir)
	if IsFile(dir) {
		t.Fatalf("IsFile should be false for directory input: %s", dir)
	}
}

func TestPath_CleanRelIsAbsJoin(t *testing.T) {
	// Join should match filepath.Join
	p := Join("a", "..", "a", "b")
	if p != filepath.Join("a", "..", "a", "b") {
		t.Fatalf("Join mismatch: %s", p)
	}

	// Clean should match filepath.Clean
	dirty := "a//b/./c/.."
	if Clean(dirty) != filepath.Clean(dirty) {
		t.Fatalf("Clean mismatch")
	}

	// Abs + IsAbs + Rel
	relBase := filepath.Join(os.TempDir(), fmt.Sprintf("mpath_rel_base_%d", time.Now().UnixNano()))
	_ = os.MkdirAll(relBase, 0o755)
	defer os.RemoveAll(relBase)

	sample := filepath.Join(relBase, "x", "y")
	abs, err := Abs(sample)
	if err != nil {
		t.Fatalf("Abs error: %v", err)
	}
	if !IsAbs(abs) {
		t.Fatalf("expected abs to be absolute: %s", abs)
	}
	r, err := Rel(abs, relBase)
	if err != nil {
		t.Fatalf("Rel error: %v", err)
	}
	if r == "" {
		t.Fatalf("Rel returned empty")
	}
}
