package utils

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestNewFileName(t *testing.T) {
	fmt.Println(filepath.Base("/foo/bar/baz.js"))
	fmt.Println(filepath.Split("/foo/bar/baz.js"))
	fmt.Println(filepath.Ext("/foo/bar/baz.js"))
}
