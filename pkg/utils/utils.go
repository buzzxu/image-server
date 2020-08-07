package utils

import (
	"fmt"
	"github.com/buzzxu/boys/types"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/**
 */
func IfImage(buff []byte) (bool, string) {
	contentType := http.DetectContentType(buff)
	switch contentType {
	case "image/jpeg", "image/jpg", "image/webp", "image/gif", "image/png":
		return true, contentType
	default:
		return false, ""
	}
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
	if suffix == extension {
		return filename
	}
	path := strings.TrimSuffix(filename, suffix)
	return path + extension
}

func GenFileNameByType(data *[]byte) (string, error) {
	if flag, contentType := IfImage(*data); flag {
		id := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		suffix := func() string {
			switch contentType {
			case "image/webp":
				return ".webp"
			case "image/jpeg", "images/jpg":
				return ".jpg"
			case "image/png":
				return ".png"
			case "image/gig":
				return ".gif"
			}
			return ".webp"
		}
		return id + suffix(), nil
	}
	return "", types.NewError(400, "非图片禁止上传")
}
func Base64ImagePrefix(data *[]byte) (string, error) {
	if flag, contentType := IfImage(*data); flag {
		return "data:" + contentType + ";base64,", nil
	}
	return "", types.NewError(400, "非图片无法区分图片类型")
}

func MkDirExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			log.Println(fmt.Sprintf("服务器目录:[%s]已存在,无需创建", path))
			return nil
		}
		if os.IsNotExist(err) {
			log.Println(fmt.Sprintf("服务器目录:[%s]不存在,需创建", path))
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf(fmt.Sprintf("服务器目录:[%s]创建失败", path))
			}
			return nil
		}
	}
	return nil
}
