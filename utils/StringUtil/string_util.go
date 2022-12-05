package StringUtil

import (
	"fmt"
	"strings"
)

// IsEmpty
// 是否是空串
func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// Format
// format 字符串 按照 Sprintf 格式
func Format(str string, a ...any) string {
	s := fmt.Sprintf(str, a...)
	return s
}

// Slf4gFormat
// 使用{} 作为占位符号来替换后面的对象
func Slf4gFormat(str string, a ...any) string {
	for _, p := range a {
		newStr := fmt.Sprintf("%v", p)
		str = strings.Replace(str, "{}", newStr, 1)
	}
	return str
}

// ArrayToStringParam
// ArrayToString 的参数struct.
type ArrayToStringParam struct {
	// 数组指针
	arr *[]any
	// 起始. 结束. 间隔字符串
	start, end, separator string
}

// ArrayToString
// 把一个数组序列化成字符串
func ArrayToString(param *ArrayToStringParam) string {
	builder := strings.Builder{}
	if !IsEmpty(param.start) {
		builder.Grow(len(param.start))
	}

	l := len(*param.arr)
	sl := len(param.separator)
	for i, c := range *param.arr {
		s := fmt.Sprintf("%v", c)
		builder.Grow(len(s))
		builder.WriteString(s)

		if i < l-1 {
			builder.Grow(sl)
			builder.WriteString(param.separator)
		}
	}

	if !IsEmpty(param.end) {
		builder.Grow(len(param.end))
	}
	return builder.String()
}
