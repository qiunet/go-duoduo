package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetLogger_Smoke(t *testing.T) {
	t.Run("同名复用与输出冒烟", func(t *testing.T) {
		// 对齐 Java TestLogger.loggerLoop：可调用各级别输出；另验证同名复用
		a := GetLogger(string(Test))
		b := GetLogger(string(Test))
		if a != b {
			t.Fatal("same name should reuse entry")
		}

		var buf bytes.Buffer
		oldOut, oldLevel := base.Out, base.GetLevel()
		base.SetOutput(&buf)
		base.SetLevel(logrus.DebugLevel)
		t.Cleanup(func() {
			base.SetOutput(oldOut)
			base.SetLevel(oldLevel)
		})

		for i := 0; i < 5; i++ {
			msg := "我们\n是中国人abc\ndccc"
			switch {
			case i < 3:
				a.Info(msg)
			case i == 3:
				a.Error(msg)
			default:
				a.Debug(msg)
			}
		}
		out := buf.String()
		if !strings.Contains(out, "我们") || !strings.Contains(out, "是中国人") {
			t.Fatalf("output missing chinese: %s", out)
		}
	})
}
