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
	for index, filename := range upload.Files {
		mw.ReadImage(filename)
		newFileName := utils.NewFileName(filename)
		err := storeFile(newFileName, mw.GetImageBlob())
		if err != nil {
			return nil, err
		}
		_, webp := upload.Params["webp"]
		if webp {
			err = generatorWebp(filename, mw)
			if err != nil {
				return nil, err
			}
		}
		paths[index] = newFileName

		mw.Clear()
	}
	return paths, nil
}

func (image *Local) Download(download *storage.Download) ([]byte, error) {

	blob, err := readFile(download.Context, download.FileName)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(blob)

	return mw.GetImageBlob(), err
}

func (image *Local) Destory() {
	imagick.Terminate()
}

func generatorWebp(filename string, mw *imagick.MagickWand) error {
	return generatorImage(filename, "webp", mw)
}

func generatorImage(filename string, extension string, mw *imagick.MagickWand) error {
	mw.Clear()
	mw.ReadImage(filename)
	err := mw.WriteImage(utils.NewFileNameExt(filename, extension))
	if err != nil {
		return err
	}
	return nil
}

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

func storeFile(filename string, blob []byte) error {
	return ioutil.WriteFile(filename, blob, 0666)
}
