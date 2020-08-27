package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/buzzxu/boys/types"
	"github.com/satori/go.uuid"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
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

func Byte2Image(data *[]byte) (img image.Image, err error) {
	filetype := http.DetectContentType(*data)
	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewBuffer(*data))
		if err != nil {
			fmt.Println("jpeg error")
			return nil, err
		}
	case "image/gif":
		img, err = gif.Decode(bytes.NewBuffer(*data))
		if err != nil {
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(bytes.NewBuffer(*data))
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}

//解析 Hex
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func GetUrlReader(url string) (r *bytes.Reader, err error) {
	buff, err := GetUrlBuffer(url)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*buff), nil
}

func GetUrlBuffer(url string) (*[]byte, error) {
	var err error
	if url[0:4] == "file" {
		buff, err := ioutil.ReadFile(url[6:])
		if err != nil {
			return nil, err
		}
		return &buff, nil
	} else {
		var resp *http.Response
		if url[0:5] == "https" {
			c := &http.Client{
				Transport: tr,
			}
			resp, err = c.Get(url)
		} else {
			resp, err = http.Get(url)
		}
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		buff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return &buff, nil
	}
}
