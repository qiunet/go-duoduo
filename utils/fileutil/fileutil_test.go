package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileExists(t *testing.T) {
	t.Run("不存在与存在", func(t *testing.T) {
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
	})
}

func TestCreateAppendGetContent(t *testing.T) {
	t.Run("创建追加并读取UTF8", func(t *testing.T) {
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
	})
}

func TestAppendTailCopyLikeJava(t *testing.T) {
	t.Run("多次中文追加尾读复制与增量TailFileFrom", func(t *testing.T) {
		// 对齐 Java TestFileUtil.appendToFile
		dir := t.TempDir()
		path := filepath.Join(dir, "Test.txt")
		path1 := filepath.Join(dir, "Test1.txt")

		for i := 0; i < 5; i++ {
			if err := AppendToFile(path, fmt.Sprintf("中Ha%d", i)); err != nil {
				t.Fatal(err)
			}
		}
		list, err := TailFile(path, 10)
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 5 {
			t.Fatalf("tail size=%d", len(list))
		}
		for i := 0; i < len(list); i++ {
			want := fmt.Sprintf("中Ha%d", i)
			if list[i] != want {
				t.Fatalf("list[%d]=%q want %q", i, list[i], want)
			}
		}

		startPos, err := GetFileLength(path)
		if err != nil {
			t.Fatal(err)
		}

		if err := Copy(path, path1); err != nil {
			t.Fatal(err)
		}
		wantContent := "中Ha0\n中Ha1\n中Ha2\n中Ha3\n中Ha4\n"
		got, err := GetFileContent(path1)
		if err != nil || got != wantContent {
			t.Fatalf("copy content=%q err=%v", got, err)
		}

		srcContent, err := GetFileContent(path1)
		if err != nil {
			t.Fatal(err)
		}
		if err := CreateFileWithContent(path1, srcContent); err != nil {
			t.Fatal(err)
		}
		got, err = GetFileContent(path1)
		if err != nil || got != wantContent {
			t.Fatalf("after create content=%q err=%v", got, err)
		}
		if err := DeleteFile(path1); err != nil {
			t.Fatal(err)
		}

		if err := AppendToFile(path, "中Ha5"); err != nil {
			t.Fatal(err)
		}
		list, err = TailFileFrom(path, startPos, 10)
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 1 || list[0] != "中Ha5" {
			t.Fatalf("incremental tail=%v", list)
		}
		if !FileExists(path) {
			t.Fatal("file should exist")
		}
		if err := DeleteFile(path); err != nil {
			t.Fatal(err)
		}
		if FileExists(path) {
			t.Fatal("file should be deleted")
		}
	})
}

func TestCopyMove(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := WriteFile(src, "abc"); err != nil {
		t.Fatal(err)
	}

	t.Run("复制文件", func(t *testing.T) {
		if err := Copy(src, dst); err != nil {
			t.Fatal(err)
		}
		got, _ := GetFileContent(dst)
		if got != "abc" {
			t.Fatalf("copy got %q", got)
		}
	})

	t.Run("复制UTF8多行全文一致", func(t *testing.T) {
		utfSrc := filepath.Join(dir, "utf-src.txt")
		utfDst := filepath.Join(dir, "utf-dst.txt")
		content := "中Ha0\n中Ha1\n中Ha2\n"
		if err := CreateFileWithContent(utfSrc, content); err != nil {
			t.Fatal(err)
		}
		if err := Copy(utfSrc, utfDst); err != nil {
			t.Fatal(err)
		}
		got, err := GetFileContent(utfDst)
		if err != nil || got != content {
			t.Fatalf("got=%q err=%v", got, err)
		}
	})

	t.Run("移动到目录", func(t *testing.T) {
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
	})
}

func TestReadFileLinesAndTail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lines.txt")
	content := "a\nb\nc\nd\n"
	if err := WriteFile(path, content); err != nil {
		t.Fatal(err)
	}

	t.Run("按行读取", func(t *testing.T) {
		lines, err := ReadFileLines(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(lines, ",") != "a,b,c,d" {
			t.Fatalf("lines=%v", lines)
		}
	})

	t.Run("读取末尾行", func(t *testing.T) {
		tail, err := TailFile(path, 2)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(tail, ",") != "c,d" {
			t.Fatalf("tail=%v", tail)
		}
	})
}

func TestDeleteAndClean(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := WriteFile(filepath.Join(sub, "f.txt"), "1"); err != nil {
		t.Fatal(err)
	}

	t.Run("递归列出文件", func(t *testing.T) {
		files, err := ListFiles(dir, nil)
		if err != nil || len(files) != 1 {
			t.Fatalf("list=%v err=%v", files, err)
		}
	})

	t.Run("清理目录内容", func(t *testing.T) {
		if err := CleanDirectory(dir); err != nil {
			t.Fatal(err)
		}
		entries, _ := os.ReadDir(dir)
		if len(entries) != 0 {
			t.Fatalf("dir not cleaned: %v", entries)
		}
	})

	t.Run("删除目录自身", func(t *testing.T) {
		if err := WriteFile(filepath.Join(dir, "x.txt"), "x"); err != nil {
			t.Fatal(err)
		}
		if err := DeleteFile(dir); err != nil {
			t.Fatal(err)
		}
		if FileExists(dir) {
			t.Fatal("dir should be deleted")
		}
	})
}

func TestGetFileLength(t *testing.T) {
	t.Run("返回字节长度", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "len.txt")
		if err := WriteFile(path, "12345"); err != nil {
			t.Fatal(err)
		}
		n, err := GetFileLength(path)
		if err != nil || n != 5 {
			t.Fatalf("len=%d err=%v", n, err)
		}
	})
}

func TestWriteStringToFile_RejectDirectory(t *testing.T) {
	t.Run("写入目录报错", func(t *testing.T) {
		dir := t.TempDir()
		err := WriteStringToFile(dir, "x", false, "")
		if err == nil {
			t.Fatal("expected error writing to directory")
		}
	})
}

func TestChangeListener(t *testing.T) {
	t.Run("文件变动回调", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "watch.txt")
		if err := WriteFile(path, "v1"); err != nil {
			t.Fatal(err)
		}
		ch := make(chan string, 1)
		if err := ChangeListener(path, func(p string) { ch <- p }); err != nil {
			t.Fatal(err)
		}
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
	})
}
