package strutil

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"utils/mathutil"
)

const (
	// EmptyString 空字符串。对应 Java StringUtil.EMPTY_STRING。
	EmptyString = ""
)

var (
	chineseRegex = regexp.MustCompile(`^[\x{4E00}-\x{9FA5}]*$`)
	numRegex     = regexp.MustCompile(`^-?[0-9]+$`)
	randChars    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hanStart     = 0x4e00
	hanEnd       = 0x9fa5
)

// IsEmptyString 判断是否空字符串（nil 语义在 Go 中即 ""）。对应 Java StringUtil.isEmpty。
func IsEmptyString(str string) bool {
	return len(str) == 0
}

// SplitString 按分隔符精确切分（非正则），保留空段。
// 例: SplitString(",", ",") -> ["", ""]。对应 Java StringUtil.split。
func SplitString(srcStr, splitStr string) []string {
	if srcStr == "" {
		return []string{}
	}
	if splitStr == "" {
		return []string{srcStr}
	}

	ret := make([]string, 0)
	before := 0
	cursor := 0
	splitValid := false
	for i := 0; i < len(srcStr); i++ {
		if !splitValid && srcStr[i] == splitStr[0] {
			splitValid = true
			cursor = i
		}
		if splitValid {
			if srcStr[i] != splitStr[i-cursor] {
				splitValid = false
				continue
			}
			if i-cursor == len(splitStr)-1 {
				ret = append(ret, srcStr[before:cursor])
				before = i + 1
				splitValid = false
			}
		}
	}
	ret = append(ret, srcStr[before:])
	return ret
}

// ConversionInt 将字符串切片转为 int 切片。对应 Java StringUtil.conversion(..., Integer.class)。
func ConversionInt(parts []string) ([]int, error) {
	out := make([]int, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

// ConversionIntSplit 先切分再转为 int。对应 Java StringUtil.conversion(src, split, Integer.class)。
func ConversionIntSplit(src, split string) ([]int, error) {
	return ConversionInt(SplitString(src, split))
}

// ArraysToString 用分隔符拼接。对应 Java StringUtil.arraysToString(arrays, separator)。
func ArraysToString[T any](arrays []T, separator string) string {
	return ArraysToStringBound(arrays, "", "", separator)
}

// ArraysToStringBound 带起止符拼接。对应 Java StringUtil.arraysToString(arrays, start, end, separator)。
func ArraysToStringBound[T any](arrays []T, start, end, separator string) string {
	if len(arrays) == 0 {
		return start + end
	}
	return ArraysToStringRange(arrays, start, end, 0, len(arrays)-1, separator)
}

// ArraysToStringRange 按下标区间拼接。对应 Java StringUtil.arraysToString(..., startIndex, endIndex, ...)。
func ArraysToStringRange[T any](arrays []T, start, end string, startIndex, endIndex int, separator string) string {
	var b strings.Builder
	b.WriteString(start)
	for i := startIndex; i <= endIndex && i < len(arrays); i++ {
		if i > startIndex {
			b.WriteString(separator)
		}
		b.WriteString(objectToString(arrays[i]))
	}
	b.WriteString(end)
	return b.String()
}

// MapToString 将 map 拼成字符串。对应 Java StringUtil.mapToString。
func MapToString[K comparable, V any](m map[K]V, keyValSeparator, separator, start, end string) string {
	var b strings.Builder
	b.WriteString(start)
	first := true
	for k, v := range m {
		if !first {
			b.WriteString(separator)
		}
		first = false
		b.WriteString(fmt.Sprint(k))
		b.WriteString(keyValSeparator)
		b.WriteString(fmt.Sprint(v))
	}
	b.WriteString(end)
	return b.String()
}

// MapToStringSimple map 拼串（无起止符）。对应 Java StringUtil.mapToString(map, kvSep, sep)。
func MapToStringSimple[K comparable, V any](m map[K]V, keyValSeparator, separator string) string {
	return MapToString(m, keyValSeparator, separator, EmptyString, EmptyString)
}

// objectToString 对象转字符串；切片/数组对齐 Java Arrays.toString 风格。
func objectToString(obj any) string {
	if obj == nil {
		return "<nil>"
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		var b strings.Builder
		b.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprint(v.Index(i).Interface()))
		}
		b.WriteByte(']')
		return b.String()
	default:
		return fmt.Sprint(obj)
	}
}

// Regex 是否完全匹配正则。对应 Java StringUtil.regex（Matcher.matches）。
func Regex(matchStr, regex string) bool {
	if IsEmptyString(matchStr) {
		return false
	}
	re, err := regexp.Compile("^(?:" + regex + ")$")
	if err != nil {
		return false
	}
	return re.MatchString(matchStr)
}

// IsNum 是否为整数数字串（可带负号）。对应 Java StringUtil.isNum。
func IsNum(numStr string) bool {
	if numStr == "" {
		return false
	}
	return numRegex.MatchString(numStr)
}

// RegexChinese 是否全是汉字。对应 Java StringUtil.regexChinese。
func RegexChinese(matchStr string) bool {
	if IsEmptyString(matchStr) {
		return false
	}
	return chineseRegex.MatchString(matchStr)
}

// Format 按 {0}{1} 占位格式化。对应 Java StringUtil.format。
func Format(pattern string, params ...any) string {
	out := pattern
	for i, p := range params {
		ph := fmt.Sprintf("{%d}", i)
		out = strings.ReplaceAll(out, ph, fmt.Sprint(p))
	}
	return out
}

// Slf4jFormat 按 {} 顺序占位格式化。对应 Java StringUtil.slf4jFormat。
func Slf4jFormat(pattern string, params ...any) string {
	var b strings.Builder
	pi := 0
	for i := 0; i < len(pattern); {
		if i+1 < len(pattern) && pattern[i] == '{' && pattern[i+1] == '}' {
			if pi < len(params) {
				b.WriteString(fmt.Sprint(params[pi]))
				pi++
			} else {
				b.WriteString("{}")
			}
			i += 2
			continue
		}
		b.WriteByte(pattern[i])
		i++
	}
	return b.String()
}

// GetIntHexVal int 转 8 位大写十六进制。对应 Java StringUtil.getIntHexVal。
func GetIntHexVal(val int) string {
	return fmt.Sprintf("%08X", uint32(val))
}

// GetByteHexVal byte 转 2 位大写十六进制。对应 Java StringUtil.getByteHexVal。
func GetByteHexVal(val byte) string {
	return fmt.Sprintf("%02X", val)
}

// GetShortHexVal short 转 4 位大写十六进制。对应 Java StringUtil.getShortHexVal。
func GetShortHexVal(val int16) string {
	return fmt.Sprintf("%04X", uint16(val))
}

// RandomString 生成指定长度随机串。对应 Java StringUtil.randomString。
func RandomString(count int) string {
	if count <= 0 {
		return EmptyString
	}
	var b strings.Builder
	b.Grow(count)
	for i := 0; i < count; i++ {
		b.WriteByte(randChars[mathutil.RandomIntN(len(randChars))])
	}
	return b.String()
}

// Repeated 将字符串重复 count 次。对应 Java StringUtil.repeated。
func Repeated(s string, count int) string {
	if count <= 0 {
		return EmptyString
	}
	return strings.Repeat(s, count)
}

// PowerfulTrim 裁剪两端非「可打印 ASCII / 汉字」字符。对应 Java StringUtil.powerfulTrim。
func PowerfulTrim(str string) string {
	if str == "" {
		return str
	}
	runes := []rune(str)
	start, length := 0, len(runes)
	for start < length && !isAllowRune(runes[start]) {
		start++
	}
	for length > start && !isAllowRune(runes[length-1]) {
		length--
	}
	if start > 0 || length < len(runes) {
		return string(runes[start:length])
	}
	return str
}

// isAllowRune 判断是否允许保留的字符。
func isAllowRune(ch rune) bool {
	if ch >= 33 && ch <= 126 {
		return true
	}
	return ch >= rune(hanStart) && ch <= rune(hanEnd)
}

// GetMixedStringLength 中英混合长度：汉字计 2，其它计 1。对应 Java StringUtil.getMixedStringLength。
func GetMixedStringLength(input string) int {
	length := 0
	for _, c := range input {
		if c >= rune(hanStart) && c <= rune(hanEnd) {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// EqualsString 比较两个字符串（含都为空）。对应 Java StringUtil.equals。
func EqualsString(str1, str2 string) bool {
	return str1 == str2
}
