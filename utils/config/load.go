package config

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Load 从 YAML 文件加载配置到 dst（须为指向 struct 的非 nil 指针）。
// 覆盖顺序：default tag → YAML → env tag。
func Load(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: read %s: %w", path, err)
	}
	return LoadBytes(data, dst)
}

// LoadBytes 从 YAML 字节加载配置到 dst（须为指向 struct 的非 nil 指针）。
// 覆盖顺序：default tag → YAML → env tag。
func LoadBytes(data []byte, dst any) error {
	v, err := structPtrValue(dst)
	if err != nil {
		return err
	}

	root := map[string]any{}
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("config: yaml unmarshal: %w", err)
		}
	}

	if err := applyDefaults(v); err != nil {
		return err
	}
	if err := bindYAML(v, root); err != nil {
		return err
	}
	if err := applyEnv(v); err != nil {
		return err
	}
	return nil
}

// structPtrValue 校验 dst 为指向 struct 的非 nil 指针并返回可设置的 Value。
func structPtrValue(dst any) (reflect.Value, error) {
	if dst == nil {
		return reflect.Value{}, fmt.Errorf("config: dst is nil")
	}
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return reflect.Value{}, fmt.Errorf("config: dst must be non-nil pointer to struct")
	}
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("config: dst must point to struct, got %s", elem.Kind())
	}
	return elem, nil
}
