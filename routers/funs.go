package routers

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/buzzxu/boys/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"image-server/pkg/conf"
	"image-server/pkg/knife"
	"image-server/pkg/utils"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

//裁剪
func crop(c echo.Context) error {
	params := c.FormValue("params")
	if params == "" {
		return E(c, types.NewHttpError(http.StatusBadRequest, "params is nil"))
	}
	var blob *[]byte
	url := c.FormValue("url")
	if url != "" {
		var err error
		blob, err = utils.GetUrlBuffer(url)
		if err != nil {
			return jsonErrorHandler(err, c)
		}
	} else if c.FormValue("base64") != "" {
		data := c.FormValue("base64")
		image := data
		index := strings.Index(data, ",")
		if index > 0 {
			image = data[index+1:]
		}
		bytes, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			return jsonErrorHandler(err, c)
		}
		blob = &bytes
	} else {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		files, exists := form.File["file"]
		if !exists {
			return E(c, types.NewHttpError(http.StatusBadRequest, "file or base64 is nil"))
		}
		file := files[0]
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
			blob = &buff
		} else {
			return E(c, types.NewHttpError(http.StatusBadRequest, fmt.Sprintf("%s不是图片,服务器拒绝上传", file.Filename)))
		}
	}
	crops := &[]knife.CropParam{}
	jsoniter.Unmarshal([]byte(params), crops)
	datas, err := knife.Crop(blob, crops)
	if err != nil {
		return jsonErrorHandler(err, c)
	}
	return c.JSON(200, datas)
}

//合并
func composite(c echo.Context) error {
	var composites []knife.CompositeParam
	if err := c.Bind(&composites); err != nil {
		return jsonErrorHandler(err, c)
	}
	blob, err := knife.Composite(composites)
	if err != nil {
		return jsonErrorHandler(err, c)
	}
	return c.Blob(200, "image/png", *blob)
}
