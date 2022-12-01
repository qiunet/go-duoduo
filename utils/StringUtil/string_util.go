package StringUtil

import (
	fmt "fmt"
	"strings"
)

// IsEmpty
// 是否是空串
func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// Format
// format 字符串
func Format(f string, a ...any) string {
	s := fmt.Sprintf(f, a)
	return s
}
