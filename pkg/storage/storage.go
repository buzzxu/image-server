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

		Download(download *Download) ([]byte, error)

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
		Files  []string
		Params map[string][]string
	}
	//读取的参数
	Download struct {
		Params   map[string]string
		Blod     []byte
		Folders  []*string
		FileName string
		FileExt  string
		Context  context.Context
	}
)

var Storager Storage

var NewStorage func(t string) Storage

func init() {
	Storager = NewStorage(conf.Config.Type)
}

func Register(ns func(t string) Storage) {
	NewStorage = ns
}
