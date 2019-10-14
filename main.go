package main

import (
	"fmt"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/storage/aliyun"
	"image-server/pkg/storage/local"
	"image-server/routers"
	"runtime"
)

func init() {
	storage.Register(func(t string) storage.Storage {
		switch conf.Config.Type {
		case "local":
			return &local.Local{}
		case "aliyun":
			return &aliyun.Aliyun{}
		}
		return &local.Local{}
	})
}

func main() {

	runtime.GOMAXPROCS(conf.Config.MaxProc)
	storage.Storager.Init()
	defer storage.Storager.Destory()

	echo := routers.New()
	echo.Logger.Fatal(echo.Start(fmt.Sprintf(":%d", conf.Config.Port)))
}
