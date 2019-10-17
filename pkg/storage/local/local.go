package local

import (
	"context"
	"github.com/buzzxu/boys/common/strs"
	"github.com/buzzxu/boys/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/imagemagick"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Local struct {
	root string
}

const (
	key_img_default = "img:default"
	key_prefix      = "notfound:"
)

func (image *Local) Init() {
	imagick.Initialize()
	redisConnect()
	//加载默认图片到redis
	loadDefaultImg()
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
	numfiles := len(upload.Blobs)
	paths := make([]string, numfiles)
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	/*ch := make(chan string, len(upload.FileNames))
	for index, blob := range upload.Files {
		//上传图片到本地硬盘
		go uploadToLocalHard(upload.FileNames[index], blob, upload, mw, ch)
	}
	for i := 0; i < numfiles; i++ {
		paths[i] = <-ch
	}*/
	for index := 0; index < numfiles; index++ {
		webp, exist := upload.Params["webp"]
		newFileName := utils.NewFileName(upload.Folder, upload.Keys[index])
		if exist {
			webpPath, err := generatorImage(upload.Blobs[index], newFileName, ".webp", upload.Resize, mw)
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
		if err = mwStoreFile(newFileName, upload.Resize, upload.Blobs[index], mw); err != nil {
			return nil, err
		}

		if upload.Thumbnail != "" {
			generatorThumbnailImage(upload.Blobs[index], newFileName, upload.Thumbnail, mw)
		}
		paths[index] = newFileName
	}
	return paths, nil
}

func (image *Local) Download(download *storage.Download) (*[]byte, string, error) {
	var (
		blob *[]byte
		err  error
	)
	//从本地硬盘读取图片
	blob, err = loadImageFromHardDrive(download)
	return blob, http.DetectContentType(*blob), err
}

func (image *Local) Delete(del *storage.Delete) (bool, error) {
	var wg sync.WaitGroup
	wg.Add(len(del.Keys))
	for _, key := range del.Keys {
		go delLocalHard(key, del.Context, del.Logger, &wg)
	}
	wg.Wait()
	return true, nil
}
func (image *Local) Destory() {
	imagick.Terminate()
	cache.Close()
}

func getDefaultImag() *[]byte {
	blob, err := cache.Get(key_img_default).Bytes()
	if err != nil {
		loadDefaultImg()
		blob, err = cache.Get(key_img_default).Bytes()
	}
	return &blob
}

//删除图片
func delLocalHard(file string, context context.Context, logger echo.Logger, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	//验证是否有此图片
	if blob, err := readFile(context, file); err != nil {
		logger.Errorf("文件:%s,删除失败,原因:无法获取图片信息", file)
	} else {
		if utils.IfImage(blob) {
			suffix := filepath.Ext(file)
			path := strings.TrimSuffix(file, suffix)
			files, err := filepath.Glob(filepath.Join(conf.Config.Storage, path) + "*")
			if err != nil {
				logger.Errorf("路径:%s,查找文件失败,原因:%s", path, err.Error())
				return
			}
			for _, f := range files {
				if err := os.Remove(f); err != nil {
					logger.Errorf("文件:%s,删除失败,原因:%s", f, err.Error()) /**/
					continue
				}
			}
		} else {
			logger.Warnf("文件:%s,删除失败,原因:非图片无法删除", file)
		}
	}
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

func loadDefaultImg() {
	if blob, error := ioutil.ReadFile(filepath.Join(conf.Config.Storage, conf.Config.DefaultImg)); error == nil {
		cache.Set(key_img_default, blob, 0)
	} else {
		log.Fatalf("读取默认图片[%s]失败,无法缓存图片", conf.Config.DefaultImg)
	}
}

//获取图片
func loadImageFromHardDrive(download *storage.Download) (*[]byte, error) {
	var (
		blob []byte
		err  error
		mw   *imagick.MagickWand
	)
	key := strs.HashSHA1(download.URL)
	keyNotfound := key_prefix + key
	//如果不存在此图像 直接返回404
	if cache.Exists(keyNotfound).Val() > 0 {
		return getDefaultImag(), types.ErrNotFound
	}
	//从缓存中获取图像
	blob, err = cache.Get(key).Bytes()
	if err != nil {
		if download.WebP || download.Format == "webp" {
			blob, err = readFileWebp(download.Context, download.Path)
		}
		if blob == nil {
			// read image from local hard driver
			blob, err = readFile(download.Context, download.Path)
			if err != nil {
				cache.Set(keyNotfound, byte('0'), conf.Config.Redis.Expiration)
				return getDefaultImag(), types.ErrNotFound
			}
		}
		if blob != nil && !download.HasParams {
			if err = cache.Set(key, blob, conf.Config.Redis.Expiration).Err(); err != nil {
				return &blob, err
			}
			return getDefaultImag(), err
		}
		mw = imagick.NewMagickWand()
		defer mw.Destroy()
		if blob == nil {
			if err = mw.ReadImage(filepath.Join(conf.Config.DefaultImg, download.Path)); err != nil {
				cache.Set(keyNotfound, byte('0'), conf.Config.Redis.Expiration)
				return getDefaultImag(), types.ErrNotFound
			}
		} else {
			mw.ReadImageBlob(blob)
		}
		//质量
		//默认75
		mw.SetCompressionQuality(75)
		if download.Quality != "" {
			quality, err := strconv.ParseUint(download.Quality, 10, 64)
			if err != nil {
				return getDefaultImag(), err
			}
			if quality < 100 {
				mw.SetCompressionQuality(uint(quality))
			}
		}

		// 缩放
		if err = imagemagick.Resize(mw, download.Thumbnail); err != nil {
			return getDefaultImag(), err
		}
		//缩略图
		if err = imagemagick.Thumbnail(mw, download.Thumbnail); err != nil {
			return getDefaultImag(), err
		}
		//格式转换
		if download.Format != "" && download.Format != "webp" {
			err = mw.SetImageFormat(download.Format)
			if err != nil {
				return getDefaultImag(), err
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
		if err = cache.Set(key, blob, conf.Config.Redis.Expiration).Err(); err != nil {
			return getDefaultImag(), err
		}
	}
	return &blob, err
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
