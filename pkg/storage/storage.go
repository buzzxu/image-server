package storage

import (
	"context"
	"image-server/pkg/conf"
)

type (
	//图片存储接口
	Storage interface {
		Init()

		Check(params map[string]string)

		Upload(upload *Upload) ([]string, error)

		Download(download *Download) ([]byte, string, error)

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
		Files     []*[]byte
		FileNames []string
		Folder    string
		Params    map[string][]string
	}
	//读取的参数
	Download struct {
		Params   map[string]string
		Blod     []byte
		Folder   string
		FileName string
		Context  context.Context
		Size     string
		Format   string
		Line     bool
		WebP     bool
		Quality  string
	}
)

var Storager Storage

var NewStorage func(t string) Storage

func Register(ns func(t string) Storage) {
	println(ns)
	Storager = ns(conf.Config.Type)
}
