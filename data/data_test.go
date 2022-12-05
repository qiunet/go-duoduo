package data

import "testing"
import "github.com/qiunet/go-duoduo/utils/StringUtil"

func TestStringEmpty(t *testing.T) {
	if !StringUtil.IsEmpty("") {
		t.Fail()
	}
}
