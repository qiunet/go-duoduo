package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists 判断路径是否存在（文件或目录）。
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Move 将源文件移动到目标目录下（保留原文件名）。destDir 为目录路径。
// 对应 Java FileUtil.move。
func Move(srcPath, destDir string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("utils: move source is a directory: %s", srcPath)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	return os.Rename(srcPath, filepath.Join(destDir, filepath.Base(srcPath)))
}

// Copy 复制文件到 newPath（覆盖已存在目标）。源必须存在且为普通文件。
// 对应 Java FileUtil.copy。
func Copy(oldPath, newPath string) error {
	info, err := os.Stat(oldPath)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("utils: file [%s] is not exist or is not a file", oldPath)
	}
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return fmt.Errorf("utils: file copy: %w", err)
	}
	if dir := filepath.Dir(newPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(newPath, data, info.Mode().Perm())
}

// AppendToFile 以 UTF-8 追加内容到文件末尾，并自动追加换行符 "\n"。
// 对应 Java FileUtil.appendToFile。
func AppendToFile(path, msg string) error {
	return WriteStringToFile(path, msg, true, "\n")
}

// CreateFileWithContent 用 content 创建（或覆盖）文件，UTF-8，不追加后缀字符。
// 对应 Java FileUtil.createFileWithContent。
func CreateFileWithContent(path, content string) error {
	return WriteStringToFile(path, content, false, "")
}

// WriteFile 以 UTF-8 覆盖写入字符串，等价于 CreateFileWithContent。
func WriteFile(path, data string) error {
	return CreateFileWithContent(path, data)
}

// WriteStringToFile 以 UTF-8 写入字符串到文件。
// append 为 true 时追加；否则覆盖。endChar 拼在 data 之后再写入。
// 若父目录不存在会尝试创建。对应 Java FileUtil.writeStringToFile（charset 固定 UTF-8）。
func WriteStringToFile(path, data string, appendMode bool, endChar string) error {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return fmt.Errorf("utils: file '%s' exists but is a directory", path)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("utils: directory '%s' could not be created: %w", dir, err)
		}
	}

	flag := os.O_CREATE | os.O_WRONLY
	if appendMode {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data + endChar)
	return err
}

// TailFile 读取文件最后 lastNum 行（UTF-8）。对应 Java FileUtil.tailFile(file, lastNum)。
func TailFile(path string, lastNum int) ([]string, error) {
	return TailFileFrom(path, 0, lastNum)
}

// TailFileFrom 从 startPos 起向后的内容中，取最后 lastNum 行（UTF-8）。
// startPos>0 时按 Java 实现减 1，常用于配合 GetFileLength 做增量尾读。
// 对应 Java FileUtil.tailFile(file, startPos, lastNum)。
func TailFileFrom(path string, startPos int64, lastNum int) ([]string, error) {
	if lastNum <= 0 {
		return nil, nil
	}
	if startPos > 0 {
		startPos--
	}

	length, err := GetFileLength(path)
	if err != nil {
		return nil, err
	}
	if length <= 0 || startPos < 0 {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(startPos, io.SeekStart); err != nil {
		return nil, err
	}

	var lines []string
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(lines) > lastNum {
		lines = lines[len(lines)-lastNum:]
	}
	return lines, nil
}

// GetFileLength 返回文件长度（字节）。目录或不存在则报错；无法打开可读时返回 0。
// 对应 Java FileUtil.getFileLength。
func GetFileLength(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("utils: file is empty or null: %w", err)
	}
	if info.IsDir() {
		return 0, errors.New("utils: file is empty or null")
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, nil
	}
	_ = f.Close()
	return info.Size(), nil
}

// DeleteFile 删除文件；若为目录则递归删除其子项及自身。
// 对应 Java FileUtil.deleteFile。
func DeleteFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if err := DeleteFile(filepath.Join(path, e.Name())); err != nil {
				return err
			}
		}
	}
	return os.Remove(path)
}

// GetFileContent 以 UTF-8 读取整个文件内容。对应 Java FileUtil.getFileContent。
func GetFileContent(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("utils: file %s is not file or empty", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ListFiles 递归收集目录下满足 pred 的文件路径。
// pred 为 nil 时收集全部普通文件。对应 Java FileUtil.listFile。
func ListFiles(dir string, pred func(path string, info os.FileInfo) bool) ([]string, error) {
	var out []string
	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return out, nil
	}
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if fi.IsDir() {
			return nil
		}
		if pred == nil || pred(path, fi) {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

// CleanDirectory 删除目录内所有子文件/子目录，保留目录本身。
// 对应 Java FileUtil.cleanDirectory。
func CleanDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var first error
	for _, e := range entries {
		if err := DeleteFile(filepath.Join(dir, e.Name())); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// ReadFileLines 按行读取文件全部内容（UTF-8）。对应 Java FileUtil.readFileLines(File)。
func ReadFileLines(path string) ([]string, error) {
	var lines []string
	err := ReadFileLinesFunc(path, func(line string) {
		lines = append(lines, line)
	})
	return lines, err
}

// ReadFileLinesFunc 按行读取文件，并对每一行调用 fn。对应 Java FileUtil.readFileLines(File, Consumer)。
func ReadFileLinesFunc(path string, fn func(line string)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		fn(sc.Text())
	}
	return sc.Err()
}
