package strutil

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	toStringNullText              = "<null>"
	toStringContentStart          = "["
	toStringContentEnd            = "]"
	toStringFieldNameValueSep     = " = "
	toStringFieldSeparator        = ", "
	toStringArrayStart            = "{"
	toStringArraySeparator        = ", "
	toStringArrayEnd              = "}"
)

// ObjToString 对对象做可读字符串输出。对应 Java ToString.toString。
// 结构体字段可通过 tag `tostring:"-"` 忽略（对齐 IgnoreToString）。
func ObjToString(obj any) string {
	if obj == nil {
		return ""
	}
	return objToStringValue(reflect.ValueOf(obj))
}

// objToStringValue 将 reflect.Value 转为字符串。
func objToStringValue(v reflect.Value) string {
	if !v.IsValid() {
		return toStringNullText
	}
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return toStringNullText
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Interface())
	case reflect.Array, reflect.Slice:
		return arrToString(v)
	case reflect.Map:
		return mapToStringValue(v)
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return fmt.Sprint(t.UnixMilli())
		}
		return objDetail(v)
	default:
		return fmt.Sprint(v.Interface())
	}
}

// arrToString 数组/切片转字符串。
func arrToString(v reflect.Value) string {
	var b strings.Builder
	b.WriteString(v.Type().Name())
	if b.Len() == 0 {
		b.WriteString(v.Type().String())
	}
	b.WriteString(toStringArrayStart)
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			b.WriteString(toStringArraySeparator)
		}
		b.WriteString(objToStringValue(v.Index(i)))
	}
	b.WriteString(toStringArrayEnd)
	return b.String()
}

// mapToStringValue map 转字符串。
func mapToStringValue(v reflect.Value) string {
	var b strings.Builder
	b.WriteString(toStringArrayStart)
	iter := v.MapRange()
	first := true
	for iter.Next() {
		if !first {
			b.WriteString(toStringArraySeparator)
		}
		first = false
		b.WriteString(objToStringValue(iter.Key()))
		b.WriteString(toStringFieldNameValueSep)
		b.WriteString(objToStringValue(iter.Value()))
	}
	b.WriteString(toStringArrayEnd)
	return b.String()
}

// objDetail 输出结构体字段详情。
func objDetail(v reflect.Value) string {
	t := v.Type()
	var b strings.Builder
	b.WriteString(t.Name())
	b.WriteString(toStringContentStart)
	first := true
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" { // 未导出
			continue
		}
		if f.Tag.Get("tostring") == "-" {
			continue
		}
		if !first {
			b.WriteString(toStringFieldSeparator)
		}
		first = false
		b.WriteString(f.Name)
		b.WriteString(toStringFieldNameValueSep)
		b.WriteString(objToStringValue(v.Field(i)))
	}
	b.WriteString(toStringContentEnd)
	return b.String()
}
