package imagemagick

import (
	"github.com/buzzxu/boys/types"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"strconv"
	"strings"
)

func Auto(mgw *imagick.MagickWand) {
	mgw.StripImage()
	//mgw.AutoLevelImage()
	//mgw.AutoGammaImage()
	//mgw.AutoOrientImage()
}

//水印
func WaterMark(mgw *imagick.MagickWand) {
	if conf.Config.WaterMark.Enable {
		waterMark := conf.Config.WaterMark
		dw := imagick.NewDrawingWand()
		pw := imagick.NewPixelWand()
		defer dw.Destroy()
		defer pw.Destroy()
		pw.SetColor(waterMark.Color)
		dw.SetFillColor(pw)
		dw.SetTextEncoding("UTF-8")
		dw.SetFontSize(waterMark.PointSize)
		dw.SetGravity(waterMark.GravityType())
		dw.SetFont(waterMark.Font)
		dw.Annotation(0, 0, waterMark.Text)
		mgw.DrawImage(dw)
	}
}

//缩放
func Resize(mgw *imagick.MagickWand, size string) error {
	return zoom(size, mgw, func(width uint, height uint) error {
		return mgw.ResizeImage(width, height, imagick.FILTER_LANCZOS2)
	})
}

//缩略图
func Thumbnail(mgw *imagick.MagickWand, size string) error {
	return zoom(size, mgw, func(width uint, height uint) error {
		return mgw.ThumbnailImage(width, height)
	})
}

//缩放
func zoom(size string, mgw *imagick.MagickWand, call func(width uint, height uint) error) error {
	if size == "" {
		return nil
	}
	var width uint
	var height uint
	var err error
	if strings.ContainsAny(size, "*") {
		width, height, err = storage.ParserSize(size)
		if err != nil {
			return err
		}
		return call(width, height)
	} else if strings.HasSuffix(size, "%") {
		width = mgw.GetImageWidth()
		height = mgw.GetImageHeight()
		val, err := strconv.Atoi(size[0:strings.IndexAny(size, "%")])
		if err != nil {
			return types.NewError(400, "format error! ex.50%")
		}
		v := uint(100 / val)
		return call(width/v, height/v)
	}
	return nil
}
