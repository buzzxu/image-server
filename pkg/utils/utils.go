package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/**
 */
func IfImage(buff []byte) bool {
	filetype := http.DetectContentType(buff)
	switch filetype {
	case "image/jpeg", "image/jpg", "image/webp", "image/gif", "image/png":
		return true
	default:
		return false
	}
	return true
}

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func IfImage1(buff []byte) bool {
	incipitStr := string(buff)
	for magic, _ := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return true
		}
	}
	return false
}

func NewFileName(folder string, filename string) string {
	suffix := filepath.Ext(filename)
	id := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	return filepath.Join(folder, id+suffix)
}
func FileNameNewExt(filename string, extension string) string {
	suffix := filepath.Ext(filename)
	path := strings.TrimSuffix(filename, suffix)
	return path + extension
}

func MkDirExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			log.Println(fmt.Sprintf("目录:[%s]已存在,无需创建", path))
			return nil
		}
		if os.IsNotExist(err) {
			log.Println(fmt.Sprintf("目录:[%s]不已存在,需创建", path))
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Fatal(fmt.Sprintf("目录:[%s]创建失败", path))
				return err
			}
			return nil
		}
	}

}
