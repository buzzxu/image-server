package minio

import (
	"strings"
	"testing"
)

func TestPath(t *testing.T) {

	var a = "/asfdf/fdfda/fdfd"
	var b = "test/xu/a.jpg"
	println(strings.TrimPrefix(a, "/"))
	println(strings.TrimPrefix(b, "/"))
}
