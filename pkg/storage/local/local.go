package local

import (
	"context"
	"github.com/buzzxu/boys/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/imagemagick"
	"image-server/pkg/redis"
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

var url_prefix_length int

func (image *Local) Init() {
	imagick.Initialize()
	redis.RedisConnect()
	//加载默认图片到redis
	loadDefaultImg()
	url_prefix_length = len(conf.Config.Domain)
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
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	numfiles := len(upload.Blobs)
	paths := make([]string, numfiles)
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
		fileName := upload.Keys[index]
		if upload.Rename {
			fileName = utils.NewFileName(upload.Folder, fileName)
		} else {
			fileName = filepath.Join(upload.Folder, fileName)
		}

		if exist {
			webpPath, err := generatorImage(upload.Blobs[index], fileName, ".webp", upload.Resize, mw)
			if err != nil {
				return nil, err
			}
			//如果只是转换图片类型操作就不需要保存原图
			if webp[0] == "convert" {
				if conf.Config.Domain != "" {
					paths[index] = conf.Config.Domain + webpPath
				} else {
					paths[index] = webpPath
				}
				continue
			}
		}
		if err = mwStoreFile(fileName, upload.Resize, upload.Blobs[index], mw); err != nil {
			return nil, err
		}
		//生成缩略图
		if upload.Thumbnail != "" {
			generatorThumbnailImage(upload.Blobs[index], fileName, upload.Thumbnail, mw)
		}
		if conf.Config.Domain != "" {
			paths[index] = conf.Config.Domain + fileName
		} else {
			paths[index] = fileName
		}
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
	redis.Close()
}

func getDefaultImag() *[]byte {
	var ctx = context.Background()
	blob, err := redis.Client.Get(ctx, key_img_default).Bytes()
	if err != nil {
		loadDefaultImg()
		blob, err = redis.Client.Get(ctx, key_img_default).Bytes()
	}
	return &blob
}

// 删除图片
func delLocalHard(file string, context context.Context, logger echo.Logger, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	if url_prefix_length > 0 && strings.Contains(file, conf.Config.Domain) {
		file = file[url_prefix_length:]
	}
	//验证是否有此图片
	if blob, err := readFile(context, file); err != nil {
		logger.Errorf("文件:%s,删除失败,原因:无法获取图片信息", file)
	} else {
		if flag, _ := utils.IfImage(blob); flag {
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
	if blob, error := storage.GetDefaultImg(); error == nil {
		redis.Client.Set(context.Background(), key_img_default, blob, 0)
	} else {
		log.Fatalf("读取默认图片[%s]失败,无法缓存图片", conf.Config.DefaultImg)
	}
}

// 获取图片
func loadImageFromHardDrive(download *storage.Download) (*[]byte, error) {
	var (
		blob []byte
		err  error
		mw   *imagick.MagickWand
	)
	var ctx = context.Background()
	//key := strs.HashSHA1(download.URL)
	//struct hashcode
	key := download.Tag
	keyNotfound := key_prefix + key
	//如果不存在此图像 直接返回404
	if redis.Client.Exists(ctx, keyNotfound).Val() > 0 {
		return getDefaultImag(), types.ErrNotFound
	}
	//从缓存中获取图像
	blob, err = redis.Client.Get(ctx, key).Bytes()
	if blob == nil {
		if download.Format == "webp" {
			blob, err = readFileWebp(download.Context, download.Path)
		}
		if blob == nil || err != nil {
			// read image from local hard driver
			if blob, err = readFile(download.Context, download.Path); err != nil {
				redis.Client.Set(ctx, keyNotfound, byte('0'), conf.Config.Redis.Expiration)
				return getDefaultImag(), types.ErrNotFound
			}
		}
		if blob != nil && !download.HasParams && download.Format == "webp" {
			if err = redis.Client.Set(ctx, key, blob, conf.Config.Redis.Expiration).Err(); err != nil {
				return getDefaultImag(), types.ErrorOf(err)
			}
			return &blob, nil
		}
		mw = imagick.NewMagickWand()
		defer mw.Destroy()
		//读取blob
		mw.ReadImageBlob(blob)
		//if blob == nil {
		//	if err = mw.ReadImage(filepath.Join(conf.Config.DefaultImg, download.Path)); err != nil {
		//		cache.Set(keyNotfound, byte('0'), conf.Config.Redis.Expiration)
		//		return getDefaultImag(), types.ErrNotFound
		//	}
		//} else {
		//	mw.ReadImageBlob(blob)
		//}
		mw.SetInterlaceScheme(imagick.INTERLACE_PLANE)
		if download.Line {
			mw.SetInterlaceScheme(imagick.INTERLACE_LINE)
		} else if download.Interlace != "" {
			switch download.Interlace {
			case "line":
				mw.SetInterlaceScheme(imagick.INTERLACE_LINE)
				break
			case "plane":
				mw.SetInterlaceScheme(imagick.INTERLACE_LINE)
				break
			case "partition":
				mw.SetInterlaceScheme(imagick.INTERLACE_PARTITION)
				break
			case "jpeg":
				mw.SetInterlaceScheme(imagick.INTERLACE_JPEG)
				break
			case "png":
				mw.SetInterlaceScheme(imagick.INTERLACE_PNG)
				break
			default:
				mw.SetInterlaceScheme(imagick.INTERLACE_NO)
			}
		}

		//质量
		//默认75
		mw.SetImageCompressionQuality(75)
		if download.Quality != "" {
			quality, err := strconv.ParseUint(download.Quality, 10, 64)
			if err != nil {
				return getDefaultImag(), types.ErrorOf(err)
			}
			if quality < 100 {
				mw.SetImageCompressionQuality(uint(quality))
			}
		}

		if download.Gamma > 0 {
			mw.SetImageGamma(download.Gamma)
		}

		// 缩放
		if err = imagemagick.Resize(mw, download.Resize); err != nil {
			return getDefaultImag(), types.ErrorOf(err)
		}
		//缩略图
		if err = imagemagick.Thumbnail(mw, download.Thumbnail); err != nil {
			return getDefaultImag(), types.ErrorOf(err)
		}
		//格式转换
		if download.Format != "" && download.Format != "webp" {
			if err = mw.SetImageFormat(download.Format); err != nil {
				return getDefaultImag(), types.ErrorOf(err)
			}
		}
		//平滑度
		//mw.SetAntialias(download.Antialias)
		mw.StripImage()
		blob, err = mw.GetImageBlob()
		if err = redis.Client.Set(ctx, key, blob, conf.Config.Redis.Expiration).Err(); err != nil {
			return &blob, types.ErrorOf(err)
		}
	}
	return &blob, err
}
func generatorImage(blob *[]byte, fileName string, extension string, resize string, mw *imagick.MagickWand) (string, error) {
	nfs := utils.FileNameNewExt(fileName, extension)
	err := mw.ReadImageBlob(*blob)
	if err != nil {
		return "", err
	}
	err = imagemagick.Resize(mw, resize)
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

/*
*
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
