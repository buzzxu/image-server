package routers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"net/http"
)

func upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	folder := c.FormValue("folder")
	if folder == "" {
		return E(c, ErrorOf(http.StatusBadRequest, "folder is nil"))
	}
	files, exists := form.File["file"]
	if !exists {
		return E(c, ErrorOf(http.StatusBadRequest, "file is nil"))
	}
	blobs := make([]*[]byte, len(files))
	fileNames := make([]string, len(files))
	for index, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		buff := make([]byte, file.Size)
		_, err = src.Read(buff)
		if utils.IfImage(buff) {
			blobs[index] = &buff
			fileNames[index] = file.Filename
		} else {
			return E(c, ErrorOf(http.StatusBadRequest, fmt.Sprintf("%s不是图片,服务器拒绝上传", file.Filename)))
		}
	}
	paths, err := storage.Storager.Upload(&storage.Upload{blobs, fileNames, folder, form.Value})
	if err != nil {
		return E(c, ErrorOf(http.StatusInternalServerError, err))
	}
	return R(c, ResultOf(http.StatusOK, paths))
}

func del(c echo.Context) error {
	return R(c, ResultOf(http.StatusOK, ""))
}

func getImage(c echo.Context) error {
	//c.Response().Write() c.Request().URL.Path
	url := c.Request().URL.Path
	folder := url[7:len(url)]

	download := &storage.Download{
		Context:  c.Request().Context(),
		Folder:   folder,
		FileName: c.Param("filename"),
		Size:     c.Param("size"),
		Format:   c.Param("format"),
		Line:     c.Param("Line") != "",
		WebP:     c.Param("webp") != "",
		Quality:  c.Param("quality"),
	}
	blob, contentType, err := storage.Storager.Download(download)
	if err != nil {
		c.Logger().Errorf("", err)
		return err
	}
	return c.Blob(200, contentType, blob)
}
