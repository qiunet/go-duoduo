package strutil

import (
	"strings"
	"testing"
)

func TestStringUtil_IsEmptyAndEquals(t *testing.T) {
	t.Run("空串判断", func(t *testing.T) {
		if !IsEmptyString("") || IsEmptyString("a") {
			t.Fatal("IsEmptyString")
		}
	})
	t.Run("相等判断", func(t *testing.T) {
		if !EqualsString("a", "a") || EqualsString("a", "b") || !EqualsString("", "") {
			t.Fatal("EqualsString")
		}
	})
}

func TestStringUtil_Split(t *testing.T) {
	t.Run("保留空段", func(t *testing.T) {
		strs := SplitString("1,2,", ",")
		if len(strs) != 3 || strs[2] != "" {
			t.Fatalf("%v", strs)
		}
		strs = SplitString(",1344,2", ",")
		if len(strs) != 3 || strs[0] != "" {
			t.Fatalf("%v", strs)
		}
		strs = SplitString("1,2", ",")
		if len(strs) != 2 {
			t.Fatalf("%v", strs)
		}
		strs = SplitString(",1,2,", ",")
		if len(strs) != 4 {
			t.Fatalf("%v", strs)
		}
	})
	t.Run("空源与空分隔符", func(t *testing.T) {
		if len(SplitString("", ",")) != 0 {
			t.Fatal("empty src")
		}
		got := SplitString("abc", "")
		if len(got) != 1 || got[0] != "abc" {
			t.Fatalf("%v", got)
		}
	})
	t.Run("多字符分隔符", func(t *testing.T) {
		got := SplitString("a::b::c", "::")
		if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
			t.Fatalf("%v", got)
		}
	})
}

func TestStringUtil_Conversion(t *testing.T) {
	t.Run("字符串切片转int", func(t *testing.T) {
		got, err := ConversionInt([]string{"1", "2", "-3"})
		if err != nil || len(got) != 3 || got[0] != 1 || got[2] != -3 {
			t.Fatalf("got=%v err=%v", got, err)
		}
		if _, err := ConversionInt([]string{"1", "x"}); err == nil {
			t.Fatal("want error")
		}
	})
	t.Run("切分后转int", func(t *testing.T) {
		got, err := ConversionIntSplit("10,20,30", ",")
		if err != nil || len(got) != 3 || got[1] != 20 {
			t.Fatalf("got=%v err=%v", got, err)
		}
	})
}

func TestStringUtil_PowerfulTrim(t *testing.T) {
	t.Run("裁剪两端空白与特殊字符", func(t *testing.T) {
		// 对齐 Java TestStringUtil.testPowerfullTrim
		if PowerfulTrim("s\ts ") != "s\ts" {
			t.Fatal("1")
		}
		if PowerfulTrim("\ts\ts\t") != "s\ts" {
			t.Fatal("2")
		}
		if PowerfulTrim("abc") != "abc" {
			t.Fatal("3")
		}
		if PowerfulTrim("ㅤ") != "" {
			t.Fatal("hangul filler")
		}
		if PowerfulTrim("　尹漂亮　 　") != "尹漂亮" {
			t.Fatal("4")
		}
		if PowerfulTrim("　　﹏ℳ๓二ৡৢﺴ") != "二" {
			t.Fatal("keep han only")
		}
		if PowerfulTrim("　　ۣۣ ۣۣ ۣۣ ۣۣ ۣۣ ۣۣ ") != "" {
			t.Fatal("all disallowed")
		}
		if PowerfulTrim("") != "" {
			t.Fatal("empty")
		}
	})
}

func TestStringUtil_IsNumAndChinese(t *testing.T) {
	t.Run("数字判断", func(t *testing.T) {
		if !IsNum("123") || !IsNum("-123") || IsNum("-1ff") || IsNum("") {
			t.Fatal("isNum")
		}
	})
	t.Run("汉字判断", func(t *testing.T) {
		// 对齐 Java TestStringUtil.testChineseWord
		if !RegexChinese("我们") || RegexChinese("ab") || RegexChinese("126") || RegexChinese("^%^%#$") || RegexChinese("") {
			t.Fatal("chinese")
		}
	})
	t.Run("正则完全匹配", func(t *testing.T) {
		if !Regex("abc", "a.c") || Regex("abcd", "a.c") || Regex("", ".*") {
			t.Fatal("regex")
		}
		if Regex("a", "[") {
			t.Fatal("invalid regex should false")
		}
	})
}

func TestStringUtil_ArraysAndMap(t *testing.T) {
	t.Run("数组拼串", func(t *testing.T) {
		// 对齐 Java TestStringUtil.testArrayToString
		arr := []string{"0", "1", "2", "3", "4", "5"}
		if ArraysToString(arr, ",") != "0,1,2,3,4,5" {
			t.Fatal("plain")
		}
		if ArraysToStringBound(arr, "[", "]", ",") != "[0,1,2,3,4,5]" {
			t.Fatal("bound")
		}
		if ArraysToStringRange(arr, "[", "]", 1, 4, ",") != "[1,2,3,4]" {
			t.Fatal("range")
		}
		nested := []any{"obj1", "obj2", arr}
		if ArraysToString(nested, ",") != "obj1,obj2,[0, 1, 2, 3, 4, 5]" {
			t.Fatalf("nested=%s", ArraysToString(nested, ","))
		}
		if ArraysToStringBound([]string{}, "[", "]", ",") != "[]" {
			t.Fatal("empty")
		}
	})
	t.Run("map拼串", func(t *testing.T) {
		m := map[string]int{"a": 1}
		got := MapToString(m, "=", ",", "[", "]")
		if got != "[a=1]" {
			t.Fatalf("MapToString=%s", got)
		}
		got = MapToStringSimple(m, ":", ";")
		if got != "a:1" {
			t.Fatalf("MapToStringSimple=%s", got)
		}
	})
}

func TestStringUtil_Format(t *testing.T) {
	t.Run("编号占位", func(t *testing.T) {
		if Format("{0}叫{1}", "我", "Qiunet") != "我叫Qiunet" {
			t.Fatal("format")
		}
	})
	t.Run("slf4j占位", func(t *testing.T) {
		if Slf4jFormat("xx{}xx{}xx", "a", "b") != "xxaxxbxx" {
			t.Fatal("slf4j")
		}
		if Slf4jFormat("x{}y{}", "a") != "xay{}" {
			t.Fatal("slf4j short params")
		}
	})
}

func TestStringUtil_HexAndRandom(t *testing.T) {
	t.Run("十六进制", func(t *testing.T) {
		if GetIntHexVal(255) != "000000FF" {
			t.Fatal("int hex")
		}
		if GetByteHexVal(15) != "0F" {
			t.Fatal("byte hex")
		}
		if GetShortHexVal(255) != "00FF" {
			t.Fatal("short hex")
		}
	})
	t.Run("随机串", func(t *testing.T) {
		if RandomString(0) != "" || RandomString(-1) != "" {
			t.Fatal("empty random")
		}
		s := RandomString(8)
		if len(s) != 8 {
			t.Fatalf("len=%d", len(s))
		}
		for _, c := range s {
			if !strings.ContainsRune(randChars, c) {
				t.Fatalf("unexpected char %c", c)
			}
		}
	})
	t.Run("重复串", func(t *testing.T) {
		if Repeated("ab", 3) != "ababab" {
			t.Fatal("repeat")
		}
		if Repeated("ab", 0) != "" || Repeated("ab", -1) != "" {
			t.Fatal("repeat empty")
		}
	})
}

func TestStringUtil_MixedLength(t *testing.T) {
	t.Run("混合长度", func(t *testing.T) {
		// 对齐 Java TestStringUtil.testMixedStringLength
		if GetMixedStringLength("a行b不行") != 8 {
			t.Fatal("mixed")
		}
		if GetMixedStringLength("abcdefg") != 7 {
			t.Fatal("ascii")
		}
		if GetMixedStringLength("嘻嘻嘻") != 6 {
			t.Fatal("cn")
		}
		if GetMixedStringLength("") != 0 {
			t.Fatal("empty")
		}
	})
}
