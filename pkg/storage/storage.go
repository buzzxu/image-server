package storage

import (
	"context"
	"github.com/labstack/echo/v4"
	"image-server/pkg/conf"
	"strconv"
	"strings"
)

type (
	//图片存储接口
	Storage interface {
		Init()

		Check(params map[string]string)

		Upload(upload *Upload) ([]string, error)

		Download(download *Download) ([]byte, string, error)

		Delete(del *Delete) (bool, error)

		Destory()
	}

	KV [2]string

	Fetcher interface {
		Fetch(string) string
	}

	File interface {
		Key() string

		Exist() (bool, string, error)

		Meta() (Fetcher, error)

		Append([]byte, int64, ...KV) (int64, string, error)

		Delete() (string, error)

		Bytes() ([]byte, string, error)

		SetMeta(...KV) error
	}

	//上传的参数
	Upload struct {
		Blobs     []*[]byte
		Keys      []string
		Folder    string
		Thumbnail string
		Resize    string
		Params    map[string][]string
	}
	//读取的参数
	Download struct {
		Params    map[string]string
		Blod      []byte
		Path      string
		FileName  string
		Context   context.Context
		HasParams bool
		Resize    string
		Format    string
		Line      bool
		WebP      bool
		Quality   string
		Thumbnail string
		Interlace string
	}
	Delete struct {
		Keys    []string
		Context context.Context
		Logger  echo.Logger
	}
)

var Storager Storage

var NewStorage func(t string) Storage

func Register(ns func(t string) Storage) {
	println(ns)
	Storager = ns(conf.Config.Type)
}

func (d *Download) Resize2WidthAndHeight() (uint, uint, error) {
	return ParserSize(d.Resize)
}
func (d *Download) Thumbnail2WidthAndHeight() (uint, uint, error) {
	return ParserSize(d.Thumbnail)
}

func ParserSize(size string) (uint, uint, error) {
	resize := strings.Split(size, "*")
	swidth, err := strconv.ParseUint(resize[0], 10, 64)
	if err != nil {
		return 0, 0, echo.NewHTTPError(400, "width is not int")
	}
	sheight, err := strconv.ParseUint(resize[1], 10, 64)
	if err != nil {
		return 0, 0, echo.NewHTTPError(400, "height is not int")
	}
	width := uint(swidth)
	height := uint(sheight)
	return width, height, nil
}
