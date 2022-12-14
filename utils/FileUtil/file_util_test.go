package FileUtil

import (
	"os"
	"path"
	"testing"
)

func TestReadContent(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fail()
	}
	join := path.Join(dir, "file_util.go")
	print(ReadContent(join))
}
