package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// Filename 从路径得到文件名。
// 例如 /Users/qiunet/doc/x.txt → x.txt。对应 Java DPath.fileName。
func Filename(path string) string {
	return filepath.Base(path)
}

// Dirname 获得文件夹路径。
// 若 path 本身是目录则返回 path；否则返回父目录。
// 例如 /Users/qiunet/doc/x.txt → /Users/qiunet/doc。对应 Java DPath.dirName。
func Dirname(path string) string {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return path
	}
	return filepath.Dir(path)
}

// Join 拼接路径片段。
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Split 将路径拆成目录与文件名。
func Split(path string) (dir, file string) {
	return filepath.Split(path)
}

// ListDir 列出目录下的直接子路径（不递归）。
// pred 为 nil 时返回全部子项；否则只返回 pred 为 true 的项。
// 若 path 不是目录则返回错误。对应 Java DPath.listDir。
func ListDir(dirPath string, pred func(path string, info os.FileInfo) bool) ([]string, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("utils: %s not a directory", dirPath)
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		fi, err := e.Info()
		if err != nil {
			return nil, err
		}
		p := filepath.Join(dirPath, e.Name())
		if pred != nil && !pred(p, fi) {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}
