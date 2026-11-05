package strutil

import (
	"strings"
	"testing"
)

type sampleToString struct {
	Name   string
	Age    int
	Secret string `tostring:"-"`
}

// toStringUser 对齐 Java TestToString.User 场景。
type toStringUser struct {
	Account string
	Scores  []int
	Level   string `tostring:"-"`
}

func TestObjToString(t *testing.T) {
	t.Run("map输出", func(t *testing.T) {
		// 对齐 Java TestToString：{qiunet = qiuyang}
		got := ObjToString(map[string]string{"qiunet": "qiuyang"})
		if got != "{qiunet = qiuyang}" {
			t.Fatalf("got=%s", got)
		}
	})
	t.Run("结构体字段输出并忽略tag", func(t *testing.T) {
		s := sampleToString{Name: "a", Age: 1, Secret: "x"}
		got := ObjToString(s)
		if !strings.Contains(got, "Name") || !strings.Contains(got, "Age") {
			t.Fatalf("got=%s", got)
		}
		if strings.Contains(got, "Secret") {
			t.Fatalf("secret should ignore: %s", got)
		}
	})
	t.Run("用户结构体含切片", func(t *testing.T) {
		// 对齐 Java：User[account = qiunet, scores = {1, 2, 3}]（Go 字段名导出、切片带类型名）
		got := ObjToString(toStringUser{Account: "qiunet", Scores: []int{1, 2, 3}, Level: "3423424"})
		if got != "toStringUser[Account = qiunet, Scores = []int{1, 2, 3}]" {
			t.Fatalf("got=%s", got)
		}
		if strings.Contains(got, "Level") {
			t.Fatalf("Level should ignore: %s", got)
		}
	})
	t.Run("整型切片输出", func(t *testing.T) {
		// 对齐 Java int[]{1, 2, 3}；Go 类型名为 []int
		got := ObjToString([]int{1, 2, 3})
		if got != "[]int{1, 2, 3}" {
			t.Fatalf("got=%s", got)
		}
	})
	t.Run("nil输出空串", func(t *testing.T) {
		if ObjToString(nil) != "" {
			t.Fatal("nil")
		}
	})
}
