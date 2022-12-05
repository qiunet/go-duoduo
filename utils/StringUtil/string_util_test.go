package StringUtil

import (
	"testing"
)

// IsEmpty 是否是空串
func TestIsEmpty(t *testing.T) {
	if !IsEmpty("") {
		t.Fail()
	}
	if !IsEmpty("  ") {
		t.Fail()
	}
	if IsEmpty(" s ") {
		t.Fail()
	}
}

func TestFormat(t *testing.T) {
	format := Format("%s am %s", "I", "Qiu")
	if format != "I am Qiu" {
		t.Fail()
	}
}

func TestSlf4gFormat(t *testing.T) {
	format := Slf4gFormat("{} am {}", "I", "Qiu")
	if format != "I am Qiu" {
		t.Fail()
	}
}

func TestArrayToString(t *testing.T) {
	type args struct {
		param *ArrayToStringParam
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"numbers", args{&ArrayToStringParam{arr: &[]any{1, 2, 3}, separator: ", "}}, "1, 2, 3"},
		{"strings", args{&ArrayToStringParam{arr: &[]any{"2", "d", "s"}, separator: "-"}}, "2-d-s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayToString(tt.args.param); got != tt.want {
				t.Errorf("ArrayToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
