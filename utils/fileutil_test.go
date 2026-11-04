package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if FileExists(path) {
		t.Fatal("missing file should not exist")
	}
	if err := WriteFile(path, "x"); err != nil {
		t.Fatal(err)
	}
	if !FileExists(path) {
		t.Fatal("file should exist")
	}
}

func TestCreateAppendGetContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "c.txt")
	if err := CreateFileWithContent(path, "你好"); err != nil {
		t.Fatal(err)
	}
	if err := AppendToFile(path, "世界"); err != nil {
		t.Fatal(err)
	}
	got, err := GetFileContent(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "你好世界\n" {
		t.Fatalf("got %q", got)
	}
}

func TestCopyMove(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := WriteFile(src, "abc"); err != nil {
		t.Fatal(err)
	}
	if err := Copy(src, dst); err != nil {
		t.Fatal(err)
	}
	got, _ := GetFileContent(dst)
	if got != "abc" {
		t.Fatalf("copy got %q", got)
	}

	movedDir := filepath.Join(dir, "moved")
	if err := Move(dst, movedDir); err != nil {
		t.Fatal(err)
	}
	if FileExists(dst) {
		t.Fatal("source should be moved away")
	}
	if !FileExists(filepath.Join(movedDir, "dst.txt")) {
		t.Fatal("moved file missing")
	}
}

func TestReadFileLinesAndTail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lines.txt")
	content := "a\nb\nc\nd\n"
	if err := WriteFile(path, content); err != nil {
		t.Fatal(err)
	}
	lines, err := ReadFileLines(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(lines, ",") != "a,b,c,d" {
		t.Fatalf("lines=%v", lines)
	}

	tail, err := TailFile(path, 2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(tail, ",") != "c,d" {
		t.Fatalf("tail=%v", tail)
	}
}

func TestDeleteAndClean(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := WriteFile(filepath.Join(sub, "f.txt"), "1"); err != nil {
		t.Fatal(err)
	}
	files, err := ListFiles(dir, nil)
	if err != nil || len(files) != 1 {
		t.Fatalf("list=%v err=%v", files, err)
	}
	if err := CleanDirectory(dir); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("dir not cleaned: %v", entries)
	}
	if err := WriteFile(filepath.Join(dir, "x.txt"), "x"); err != nil {
		t.Fatal(err)
	}
	if err := DeleteFile(dir); err != nil {
		t.Fatal(err)
	}
	if FileExists(dir) {
		t.Fatal("dir should be deleted")
	}
}

func TestGetFileLength(t *testing.T) {
	path := filepath.Join(t.TempDir(), "len.txt")
	if err := WriteFile(path, "12345"); err != nil {
		t.Fatal(err)
	}
	n, err := GetFileLength(path)
	if err != nil || n != 5 {
		t.Fatalf("len=%d err=%v", n, err)
	}
}

func TestWriteStringToFile_RejectDirectory(t *testing.T) {
	dir := t.TempDir()
	err := WriteStringToFile(dir, "x", false, "")
	if err == nil {
		t.Fatal("expected error writing to directory")
	}
}

func TestChangeListener(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.txt")
	if err := WriteFile(path, "v1"); err != nil {
		t.Fatal(err)
	}
	ch := make(chan string, 1)
	if err := ChangeListener(path, func(p string) { ch <- p }); err != nil {
		t.Fatal(err)
	}
	// 直接触发一次检查，避免等 10s
	obs := changeMonitor.observers[path]
	if obs == nil {
		t.Fatal("observer not registered")
	}
	time.Sleep(10 * time.Millisecond)
	if err := WriteFile(path, "v2"); err != nil {
		t.Fatal(err)
	}
	obs.checkAndNotify()
	select {
	case p := <-ch:
		if p != path {
			t.Fatalf("path=%s", p)
		}
	case <-time.After(time.Second):
		t.Fatal("callback not fired")
	}
}
