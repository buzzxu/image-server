package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
	"image-server/pkg/conf"
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

func NewFileName(filename string) string {
	extension := filepath.Ext(filename)
	id := string(uuid.NewV4().Bytes())
	dir, _ := filepath.Split(filename)
	return filepath.Join(conf.Config.Storage, dir, id+extension)
}
func NewFileNameExt(filename string, extension string) string {
	id := string(uuid.NewV4().Bytes())
	dir, _ := filepath.Split(filename)
	return filepath.Join(conf.Config.Storage, dir, id+"."+extension)
}

func MkDirExist(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			log.Println(fmt.Sprintf("目录:[%s]已存在,无需创建", path))
			return
		}
		if os.IsNotExist(err) {
			log.Println(fmt.Sprintf("目录:[%s]不已存在,需创建", path))
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Fatal(fmt.Sprintf("目录:[%s]创建失败", path))
			} else {
				log.Println(fmt.Sprintf("目录:[%s]创建成功", path))
			}
			return
		}
	}

}
