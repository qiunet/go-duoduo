package utils

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "普通文件", path: "/Users/qiunet/doc/x.txt", want: "x.txt"},
		{name: "相对路径", path: "doc/x.txt", want: "x.txt"},
		{name: "仅文件名", path: "x.txt", want: "x.txt"},
		{name: "目录风格尾部分隔", path: filepath.Join("a", "b") + string(os.PathSeparator), want: "b"},
		{name: "根路径", path: string(os.PathSeparator), want: string(os.PathSeparator)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filename(tt.path); got != tt.want {
				t.Fatalf("Filename(%q)=%q want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestDirname(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "x.txt")
	if err := WriteFile(file, "1"); err != nil {
		t.Fatal(err)
	}

	t.Run("文件返回父目录", func(t *testing.T) {
		if got := Dirname(file); got != dir {
			t.Fatalf("Dirname(file)=%q want %q", got, dir)
		}
	})

	t.Run("目录返回自身", func(t *testing.T) {
		if got := Dirname(dir); got != dir {
			t.Fatalf("Dirname(dir)=%q want %q", got, dir)
		}
	})

	t.Run("不存在路径回退父目录", func(t *testing.T) {
		missing := filepath.Join(dir, "no-such", "a.txt")
		want := filepath.Dir(missing)
		if got := Dirname(missing); got != want {
			t.Fatalf("Dirname(missing)=%q want %q", got, want)
		}
	})
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{name: "多段", elem: []string{"a", "b", "c.txt"}, want: filepath.Join("a", "b", "c.txt")},
		{name: "单段", elem: []string{"a.txt"}, want: "a.txt"},
		{name: "空切片", elem: nil, want: ""},
		{name: "含空段", elem: []string{"a", "", "b"}, want: filepath.Join("a", "b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.elem...); got != tt.want {
				t.Fatalf("Join(%v)=%q want %q", tt.elem, got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	p := Join("a", "b", "c.txt")
	dir, file := Split(p)
	if file != "c.txt" {
		t.Fatalf("file=%q", file)
	}
	if filepath.Clean(dir) != filepath.Join("a", "b") {
		t.Fatalf("dir=%q", dir)
	}

	d2, f2 := Split("only.txt")
	if d2 != "" || f2 != "only.txt" {
		t.Fatalf("Split(only.txt)=%q %q", d2, f2)
	}
}

func TestListDir(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.log")
	sub := filepath.Join(dir, "sub")
	nested := filepath.Join(sub, "c.txt")
	if err := WriteFile(a, "1"); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(b, "2"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(nested, "3"); err != nil {
		t.Fatal(err)
	}

	t.Run("全部直接子项不递归", func(t *testing.T) {
		all, err := ListDir(dir, nil)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{a, b, sub}
		if len(all) != 3 {
			t.Fatalf("all=%v", all)
		}
		for _, w := range want {
			if !slices.Contains(all, w) {
				t.Fatalf("missing %s in %v", w, all)
			}
		}
		if slices.Contains(all, nested) {
			t.Fatalf("should not recurse into nested: %v", all)
		}
	})

	t.Run("按扩展名过滤文件", func(t *testing.T) {
		onlyTxt, err := ListDir(dir, func(path string, info os.FileInfo) bool {
			return !info.IsDir() && filepath.Ext(path) == ".txt"
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(onlyTxt) != 1 || onlyTxt[0] != a {
			t.Fatalf("onlyTxt=%v", onlyTxt)
		}
	})

	t.Run("只保留目录", func(t *testing.T) {
		onlyDirs, err := ListDir(dir, func(path string, info os.FileInfo) bool {
			return info.IsDir()
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(onlyDirs) != 1 || onlyDirs[0] != sub {
			t.Fatalf("onlyDirs=%v", onlyDirs)
		}
	})

	t.Run("空目录", func(t *testing.T) {
		empty := filepath.Join(dir, "empty")
		if err := os.Mkdir(empty, 0755); err != nil {
			t.Fatal(err)
		}
		got, err := ListDir(empty, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("empty dir got %v", got)
		}
	})

	t.Run("非目录报错", func(t *testing.T) {
		if _, err := ListDir(a, nil); err == nil {
			t.Fatal("expected error for non-directory")
		}
	})

	t.Run("不存在路径报错", func(t *testing.T) {
		_, err := ListDir(filepath.Join(dir, "missing"), nil)
		if err == nil {
			t.Fatal("expected error for missing path")
		}
		if !errors.Is(err, os.ErrNotExist) {
			// Stat 包装后通常仍可用 IsNotExist
			if !os.IsNotExist(err) {
				t.Fatalf("err=%v", err)
			}
		}
	})
}
