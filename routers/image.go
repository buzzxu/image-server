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
	/*values,err :=c.FormParams()
	if err != nil {
		return err
	}*/
	files := form.File["files"]
	fileNames := make([]string, len(files))
	for index, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		buff := make([]byte, 512)
		_, err = src.Read(buff)
		if utils.IfImage(buff) {
			fileNames[index] = file.Filename
		} else {
			return E(c, ErrorOf(http.StatusBadRequest, fmt.Sprintf("%s不是图片,服务器拒绝上传", file.Filename)))
		}
	}
	storage.Storager.Upload(&storage.Upload{fileNames, form.Value})

	return R(c, ResultOf(http.StatusOK, ""))
}

func del(c echo.Context) error {
	return R(c, ResultOf(http.StatusOK, ""))
}

func getImage(c echo.Context) error {
	//c.Response().Write()
	return c.String(http.StatusOK, c.Param("folder"))
}
