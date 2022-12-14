package FileUtil

import (
	"os"
)

// ReadContent
// 读取文件内容
func ReadContent(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
