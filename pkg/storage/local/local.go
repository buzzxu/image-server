package local

import (
	"context"
	"github.com/buzzxu/boys/types"
	"github.com/labstack/echo/v4"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/imagemagick"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	err := utils.MkDirExist(filepath.Join(conf.Config.Storage, upload.Folder))
	if err != nil {
		return nil, err
	}
	numfiles := len(upload.FileNames)
	paths := make([]string, numfiles)
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	ch := make(chan string, len(upload.FileNames))
	for index, blob := range upload.Files {
		//上传图片到本地硬盘
		go uploadToLocalHard(upload.FileNames[index], blob, upload, mw, ch)
	}
	for i := 0; i < numfiles; i++ {
		paths[i] = <-ch
	}
	return paths, nil
}

func (image *Local) Download(download *storage.Download) ([]byte, string, error) {

	var blob []byte
	var err error
	if download.WebP || download.Format == "webp" {
		blob, err = readFileWebp(download.Context, download.Path)
		if err != nil {
			return nil, "", echo.ErrNotFound
		}
	} else if blob == nil {
		// read image from local hard driver
		blob, err = readFile(download.Context, download.Path)
		if err != nil {
			return nil, "", echo.ErrNotFound
		}
	}
	if !download.HasParams {
		return blob, http.DetectContentType(blob), nil
	}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(blob)
	//质量
	//默认75
	mw.SetCompressionQuality(75)
	if download.Quality != "" {
		quality, err := strconv.ParseUint(download.Quality, 10, 64)
		if err != nil {
			return nil, "", err
		}
		if quality < 100 {
			mw.SetCompressionQuality(uint(quality))
		}
	}

	// 缩放
	if err = imagemagick.Resize(mw, download.Thumbnail); err != nil {
		return nil, "", err
	}
	//缩略图
	if err = imagemagick.Thumbnail(mw, download.Thumbnail); err != nil {
		return nil, "", err
	}
	//格式转换
	if download.Format != "" && download.Format != "webp" {
		err = mw.SetImageFormat(download.Format)
		if err != nil {
			return nil, "", err
		}
	}
	mw.SetInterlaceScheme(imagick.INTERLACE_LINE)
	if download.Interlace != "" {
		if download.Interlace == "plane" {
			mw.SetInterlaceScheme(imagick.INTERLACE_PLANE)
		}
	}

	mw.StripImage()
	blob = mw.GetImageBlob()
	return blob, http.DetectContentType(blob), err
}

func (image *Local) Destory() {
	imagick.Terminate()
}

func uploadToLocalHard(fileName string, blob *[]byte, upload *storage.Upload, mw *imagick.MagickWand, ch chan<- string) {
	webp, exist := upload.Params["webp"]
	newFileName := utils.NewFileName(upload.Folder, fileName)
	if exist {
		var webpPath string
		webpPath, _ = generatorImage(blob, newFileName, ".webp", upload.Resize, mw)
		//如果只是转换图片类型操作就不需要保存原图
		if webp[0] == "convert" {
			newFileName = webpPath
			ch <- newFileName
			return
		}
	}
	//生成缩略图
	if upload.Thumbnail != "" {
		generatorThumbnailImage(blob, newFileName, upload.Thumbnail, mw)
	}
	//保存原图
	mwStoreFile(newFileName, upload.Resize, blob, mw)
	ch <- newFileName
	return
}

func generatorImage(blob *[]byte, fileName string, extension string, resize string, mw *imagick.MagickWand) (string, error) {
	nfs := utils.FileNameNewExt(fileName, extension)
	mw.ReadImageBlob(*blob)
	err := imagemagick.Resize(mw, resize)
	if err != nil {
		return "", err
	}
	imagemagick.Auto(mw)
	err = mw.WriteImage(filepath.Join(conf.Config.Storage, nfs))
	mw.Clear()
	if err != nil {
		return "", err
	}
	return nfs, nil
}

func generatorThumbnailImage(blob *[]byte, fileName string, size string, mw *imagick.MagickWand) (string, error) {
	suffix := filepath.Ext(fileName)
	path := strings.TrimSuffix(fileName, suffix)
	nfs := path + "-thumbnail" + suffix
	mw.ReadImageBlob(*blob)
	mw.StripImage()
	err := imagemagick.Thumbnail(mw, size)
	if err != nil {
		return "", err
	}
	imagemagick.Auto(mw)
	err = mw.WriteImage(filepath.Join(conf.Config.Storage, nfs))
	mw.Clear()
	if err != nil {
		return "", err
	}
	return nfs, nil
}
func mwStoreFile(filename string, resize string, blob *[]byte, mw *imagick.MagickWand) error {
	mw.ReadImageBlob(*blob)
	err := imagemagick.Resize(mw, resize)
	if err != nil {
		return err
	}
	imagemagick.WaterMark(mw)
	imagemagick.Auto(mw)
	err = mw.WriteImage(filepath.Join(conf.Config.Storage, filename))
	mw.Clear()
	if err != nil {
		return err
	}
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
		done chan int = make(chan int) //make(chan int,runtime.GOMAXPROCS(conf.Config.MaxProc))
	)
	go func() {
		blob, err = ioutil.ReadFile(filepath.Join(conf.Config.Storage, filename))
		close(done)
	}()
	select {
	case <-ctx.Done():
		return nil, types.ErrNotFound
	case <-done:
		return blob, err
	}
}

func readFileWebp(ctx context.Context, path string) ([]byte, error) {
	webP := utils.FileNameNewExt(path, ".webp")
	return readFile(ctx, webP)
}
