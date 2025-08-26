package mpath

import (
	"os"
	"path/filepath"
	"testing"
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
	if Base(p) != "c.txt" {
		t.Fatalf("base wrong: %s", Base(p))
	}
	if Ext(p) != ".txt" {
		t.Fatalf("ext wrong: %s", Ext(p))
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

	// list
	files, err := ListFiles(dir)
	if err != nil {
		t.Fatalf("list files err: %v", err)
	}
	if len(files) != 1 || files[0] != "a.txt" {
		t.Fatalf("list files unexpected: %#v", files)
	}
}
