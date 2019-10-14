package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFileName(t *testing.T) {
	fmt.Println(filepath.Base("/foo/bar/baz.js"))
	fmt.Println(filepath.Split("/foo/bar/baz.js"))
	fmt.Println(strings.HasSuffix("50%", "%"))

}
