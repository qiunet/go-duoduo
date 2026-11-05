package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// applyDefaults 按 default tag 填充标量字段（递归嵌套 struct）。
func applyDefaults(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		fv := v.Field(i)
		if fv.Kind() == reflect.Struct && fv.Type() != reflect.TypeOf(time.Duration(0)) {
			if err := applyDefaults(fv); err != nil {
				return err
			}
			continue
		}
		def := f.Tag.Get("default")
		if def == "" {
			continue
		}
		if err := setFieldFromString(fv, def); err != nil {
			return fmt.Errorf("config: default for %s: %w", f.Name, err)
		}
	}
	return nil
}

// bindYAML 按 yaml tag（支持 a.b.c）从 map 绑定到 struct。
func bindYAML(v reflect.Value, root map[string]any) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		key := yamlTagName(f)
		if key == "-" {
			continue
		}
		fv := v.Field(i)

		if fv.Kind() == reflect.Struct && fv.Type() != reflect.TypeOf(time.Duration(0)) && !strings.Contains(key, ".") {
			sub, ok := lookupPath(root, key)
			if !ok || sub == nil {
				continue
			}
			subMap, ok := asStringKeyMap(sub)
			if !ok {
				return fmt.Errorf("config: field %s expects mapping for nested struct", f.Name)
			}
			if err := bindYAML(fv, subMap); err != nil {
				return err
			}
			continue
		}

		raw, ok := lookupPath(root, key)
		if !ok || raw == nil {
			continue
		}
		if err := setFieldFromAny(fv, raw); err != nil {
			return fmt.Errorf("config: field %s: %w", f.Name, err)
		}
	}
	return nil
}

// applyEnv 按 env tag 用环境变量覆盖叶子标量（递归嵌套 struct）。
func applyEnv(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		fv := v.Field(i)
		if fv.Kind() == reflect.Struct && fv.Type() != reflect.TypeOf(time.Duration(0)) {
			if err := applyEnv(fv); err != nil {
				return err
			}
			continue
		}
		envKey := f.Tag.Get("env")
		if envKey == "" {
			continue
		}
		val, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}
		if err := setFieldFromString(fv, val); err != nil {
			return fmt.Errorf("config: env %s for %s: %w", envKey, f.Name, err)
		}
	}
	return nil
}

// yamlTagName 取 yaml tag 名（忽略选项）；空则用字段名。
func yamlTagName(f reflect.StructField) string {
	tag := f.Tag.Get("yaml")
	if tag == "" {
		return f.Name
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return f.Name
	}
	return name
}

// lookupPath 按 a.b.c 在 map 中取值。
func lookupPath(root map[string]any, path string) (any, bool) {
	if path == "" {
		return nil, false
	}
	parts := strings.Split(path, ".")
	var cur any = root
	for _, p := range parts {
		m, ok := asStringKeyMap(cur)
		if !ok {
			return nil, false
		}
		next, ok := m[p]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

// asStringKeyMap 将 YAML 解出的 map 转为 map[string]any。
func asStringKeyMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			ks, ok := k.(string)
			if !ok {
				return nil, false
			}
			out[ks] = val
		}
		return out, true
	default:
		return nil, false
	}
}

// setFieldFromAny 将 YAML 节点值写入字段。
func setFieldFromAny(fv reflect.Value, raw any) error {
	if !fv.CanSet() {
		return fmt.Errorf("cannot set")
	}
	if fv.Type() == reflect.TypeOf(time.Duration(0)) {
		switch x := raw.(type) {
		case string:
			return setFieldFromString(fv, x)
		case int, int64, uint64, float64:
			return setFieldFromString(fv, fmt.Sprint(x)+"ns")
		default:
			return fmt.Errorf("unsupported duration value %T", raw)
		}
	}
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(fmt.Sprint(raw))
		return nil
	case reflect.Bool:
		switch x := raw.(type) {
		case bool:
			fv.SetBool(x)
			return nil
		default:
			return setFieldFromString(fv, fmt.Sprint(raw))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := toInt64(raw)
		if err != nil {
			return err
		}
		fv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := toUint64(raw)
		if err != nil {
			return err
		}
		fv.SetUint(n)
		return nil
	case reflect.Float32, reflect.Float64:
		n, err := toFloat64(raw)
		if err != nil {
			return err
		}
		fv.SetFloat(n)
		return nil
	case reflect.Slice, reflect.Map:
		// 经 yaml 再解一次，保证类型正确
		b, err := yaml.Marshal(raw)
		if err != nil {
			return err
		}
		ptr := reflect.New(fv.Type())
		if err := yaml.Unmarshal(b, ptr.Interface()); err != nil {
			return err
		}
		fv.Set(ptr.Elem())
		return nil
	default:
		return fmt.Errorf("unsupported kind %s", fv.Kind())
	}
}

// setFieldFromString 将字符串解析并写入字段。
func setFieldFromString(fv reflect.Value, s string) error {
	if !fv.CanSet() {
		return fmt.Errorf("cannot set")
	}
	if fv.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		fv.SetInt(int64(d))
		return nil
	}
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(s)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
		return nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(n)
		return nil
	default:
		return fmt.Errorf("unsupported kind %s for string value", fv.Kind())
	}
}

// toInt64 将 YAML 标量转为 int64。
func toInt64(raw any) (int64, error) {
	switch x := raw.(type) {
	case int:
		return int64(x), nil
	case int64:
		return x, nil
	case uint64:
		return int64(x), nil
	case float64:
		return int64(x), nil
	case string:
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", raw)
	}
}

// toUint64 将 YAML 标量转为 uint64。
func toUint64(raw any) (uint64, error) {
	switch x := raw.(type) {
	case int:
		if x < 0 {
			return 0, fmt.Errorf("negative")
		}
		return uint64(x), nil
	case int64:
		if x < 0 {
			return 0, fmt.Errorf("negative")
		}
		return uint64(x), nil
	case uint64:
		return x, nil
	case float64:
		if x < 0 {
			return 0, fmt.Errorf("negative")
		}
		return uint64(x), nil
	case string:
		return strconv.ParseUint(x, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to uint", raw)
	}
}

// toFloat64 将 YAML 标量转为 float64。
func toFloat64(raw any) (float64, error) {
	switch x := raw.(type) {
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(x, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float", raw)
	}
}
