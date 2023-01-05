package Resource

import (
	"fmt"
	"testing"
)

func TestResourcePath(t *testing.T) {
	path := ResourcePath("conf", "config")
	fmt.Println(path)
}
