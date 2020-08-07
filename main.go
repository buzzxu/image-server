package main

import (
	"fmt"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/storage/aliyun"
	"image-server/pkg/storage/local"
	"image-server/pkg/storage/minio"
	"image-server/pkg/storage/none"
	"image-server/pkg/storage/seaweedfs"
	"image-server/routers"
	"runtime"
)

func init() {
	storage.Register(func(t string) storage.Storage {
		switch conf.Config.Type {
		case "none":
			return &none.None{}
		case "local":
			return &local.Local{}
		case "aliyun":
			return &aliyun.Aliyun{}
		case "seaweed":
			return &seaweedfs.Seaweedfs{}
		case "minio":
			return &minio.Minio{}
		}
		return &none.None{}
	})
}

func main() {

	runtime.GOMAXPROCS(conf.Config.MaxProc)
	storage.Storager.Init()
	storage.Storager.Check(nil)
	defer storage.Storager.Destory()

	echo := routers.New()
	echo.Logger.Fatal(echo.Start(fmt.Sprintf(":%d", conf.Config.Port)))
}
