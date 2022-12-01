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
	format := Format("%d am %d", "I", "Qiu")
	if format != "I am Qiu" {
		t.Fail()
	}
}
