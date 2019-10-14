package local

import (
	"context"
	"errors"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"io/ioutil"
	"path/filepath"
)

type Local struct {
	root string
}

func (image *Local) Init() {
	imagick.Initialize()
}

func (image *Local) Check(params map[string]string) {
	//检查上传目录是否存在
	utils.MkDirExist(conf.Config.Storage)
}

func (image *Local) Upload(upload *storage.Upload) ([]string, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	paths := make([]string, len(upload.Files))
	for index, blob := range upload.Files {
		webp, exist := upload.Params["webp"]
		newFileName := utils.NewFileName(upload.Folder, upload.FileNames[index])
		if exist {
			webpPath, err := generatorImage(blob, newFileName, ".webp", mw)
			if err != nil {
				return nil, err
			}
			//如果只是转换图片类型操作就不需要保存原图
			if webp[0] == "convert" {
				newFileName = webpPath
				paths[index] = newFileName
				continue
			}
		}

		mw.ReadImageBlob(*blob)
		err := mwStoreFile(newFileName, blob, mw)
		if err != nil {
			return nil, err
		}
		paths[index] = newFileName
	}
	return paths, nil
}

func (image *Local) Download(download *storage.Download) ([]byte, string, error) {

	blob, err := readFile(download.Context, download.FileName)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(blob)

	return mw.GetImageBlob(), "", err
}

func (image *Local) Destory() {
	imagick.Terminate()
}

func generatorImage(blob *[]byte, fileName string, extension string, mw *imagick.MagickWand) (string, error) {
	mw.ReadImageBlob(*blob)
	nfs := utils.FileNameNewExt(fileName, extension)
	err := mw.WriteImage(filepath.Join(conf.Config.Storage, nfs))
	mw.Clear()
	if err != nil {
		return "", err
	}
	return nfs, nil
}

func mwStoreFile(filename string, blob *[]byte, mw *imagick.MagickWand) error {
	err := mw.WriteImage(filepath.Join(conf.Config.Storage, filename))
	mw.Clear()
	return err
}
func storeFile(filename string, blob *[]byte) error {
	return ioutil.WriteFile(filepath.Join(conf.Config.Storage, filename), *blob, 0666)
}

/**
读取文件
*/
func readFile(ctx context.Context, filename string) ([]byte, error) {
	var (
		blob []byte
		err  error
		done chan int = make(chan int)
	)
	go func() {
		blob, err = ioutil.ReadFile(filepath.Join(conf.Config.Storage, filename))
		close(done)
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("context timeout")
	case <-done:
		return blob, err
	}
}
