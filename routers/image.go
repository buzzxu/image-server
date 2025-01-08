package routers

import (
	"encoding/base64"
	"fmt"
	"github.com/buzzxu/boys/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/bytes"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	limit        int64
	cacheControl string
)

func init() {
	size, err := bytes.Parse(conf.Config.SizeLimit)
	if err != nil {
		log.Fatalf("图片限制尺寸值[%s]解析失败", conf.Config.SizeLimit)
	}
	limit = size
	cacheControl = "public,max-age=" + strconv.Itoa(conf.Config.MaxAge)
}
func upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	folder := c.FormValue("folder")
	if folder == "" {
		return E(c, types.NewHttpError(http.StatusBadRequest, "folder is nil"))
	}
	files, exists := form.File["file"]
	if !exists && c.FormValue("base64") == "" {
		return E(c, types.NewHttpError(http.StatusBadRequest, "file or base64 is nil"))
	}
	_rename := c.FormValue("rename")
	var (
		blobs     []*[]byte
		fileNames []string
	)
	var rename bool
	if _rename == "" {
		rename = true
	} else {
		rename, err = strconv.ParseBool(_rename)
		if err != nil {
			rename = true
		}
	}
	if files != nil {
		blobs = make([]*[]byte, len(files))
		fileNames = make([]string, len(files))
		for index, file := range files {
			src, err := file.Open()
			defer src.Close()
			if err != nil {
				return err
			}
			if file.Size > limit {
				return E(c, types.NewHttpError(http.StatusRequestEntityTooLarge, conf.Config.SizeLimit))
			}
			buff := make([]byte, file.Size)
			_, err = src.Read(buff)
			if flag, _ := utils.IfImage(buff); flag {
				blobs[index] = &buff
				fileNames[index] = file.Filename
			} else {
				return E(c, types.NewHttpError(http.StatusBadRequest, fmt.Sprintf("%s不是图片,服务器拒绝上传", file.Filename)))
			}
		}
	} else {
		data := c.FormValue("base64")
		image := data
		index := strings.Index(data, ",")
		if index > 0 {
			image = data[index+1:]
		}
		blob, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			return jsonErrorHandler(err, c)
		}
		blobs = make([]*[]byte, 1)
		fileNames = make([]string, 1)
		blobs[0] = &blob
		fileNames[0], _ = utils.GenFileNameByType(&blob)
		rename = false
	}

	paths, err := storage.Storager.Upload(&storage.Upload{
		Blobs:     blobs,
		Keys:      fileNames,
		Folder:    folder,
		Thumbnail: c.FormValue("thumbnail"),
		Resize:    c.FormValue("resize"),
		Params:    form.Value,
		Rename:    rename,
	})
	if err != nil {
		return jsonErrorHandler(err, c)
	}
	return R(c, paths)
}

func del(c echo.Context) error {
	files := c.QueryParams()["file"]
	if files == nil {
		return jsonErrorHandler(types.NewHttpError(http.StatusBadRequest, "请传入file参数"), c)
	}
	if _, err := storage.Storager.Delete(&storage.Delete{
		Keys:    files,
		Context: c.Request().Context(),
		Logger:  c.Logger(),
	}); err != nil {
		return jsonErrorHandler(err, c)
	}
	return RNullData(c)
}

func getImage(c echo.Context) error {
	url := c.Request().URL.Path
	path := url[7:]
	_, webp := c.QueryParams()["webp"]
	_, antialias := c.QueryParams()["antialias"]
	_, line := c.QueryParams()["line"]
	format := c.QueryParam("format")
	if format != "" {
		c.QueryParams().Del("format")
	}
	if webp {
		c.QueryParams().Del("webp")
		format = "webp"
	}
	download := &storage.Download{
		Context:   c.Request().Context(),
		Logger:    c.Logger(),
		Path:      path,
		FileName:  c.QueryParam("filename"),
		URL:       c.Request().RequestURI,
		Resize:    c.QueryParam("size"),
		Format:    format,
		Line:      line,
		Quality:   c.QueryParam("quality"),
		Thumbnail: c.QueryParam("thumbnail"),
		Interlace: c.QueryParam("interlace"),
		Antialias: antialias,
	}
	if gamma := c.QueryParam("gamma"); gamma != "" {
		if v, err := strconv.ParseFloat(gamma, 64); err == nil {
			download.Gamma = v
		}
	}
	download.HasParams = len(c.QueryParams()) > 0
	// cache
	etag := `"` + download.Etag() + `"`
	if match := c.Request().Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			return c.NoContent(http.StatusNotModified)
		}
	}
	blob, contentType, err := storage.Storager.Download(download)
	code := http.StatusOK
	if err != nil {
		if e, ok := err.(*types.Error); ok {
			code = e.Code
		} else if e, ok := err.(*echo.HTTPError); ok {
			code = e.Code
		} else {
			code = http.StatusInternalServerError
		}
		if code >= 500 {
			c.Logger().Errorf("%s 读取失败,原因:%s", path, err.Error())
		}
	} else {
		c.Response().Header().Set("Cache-Control", cacheControl)
		c.Response().Header().Set("Content-Length", strconv.Itoa(len(*blob)))
		c.Response().Header().Set("Etag", etag)
	}
	return c.Blob(code, contentType, *blob)
}
