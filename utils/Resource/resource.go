package Resource

import (
	"os"
	"path"
	"path/filepath"
)

// RootDir 得到执行文件的 root 路径
func RootDir() string {
	dir, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(dir)
	symlinks, err := filepath.EvalSymlinks(exPath)
	if err != nil {
		panic(err)
	}
	return symlinks
}

// ResourcePath 获得基于pathArray的资源路径
func ResourcePath(pathArray ...string) string {
	join := path.Join(pathArray...)
	return path.Join(RootDir(), join)
}
