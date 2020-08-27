package imagemagick

import (
	"fmt"
	"github.com/buzzxu/boys/types"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/conf"
	"image-server/pkg/objects"
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

//添加文字
func AppendText(mgw *imagick.MagickWand, x, y float64, text string, font *objects.Font) {
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()
	defer dw.Destroy()
	defer pw.Destroy()
	pw.SetColor(font.Color)
	dw.SetFillColor(pw)
	dw.SetTextEncoding("UTF-8")
	dw.SetFontSize(font.Size)
	dw.SetGravity(font.Gravity())
	dw.SetFont(font.Name)
	dw.Annotation(x, y, text)
	mgw.DrawImage(dw)
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

//裁剪
func Crop(width, height uint, x, y int, mgw *imagick.MagickWand) error {
	imCols := mgw.GetImageWidth()
	imRows := mgw.GetImageHeight()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if uint(x) >= imCols || uint(y) >= imRows {
		return fmt.Errorf("x, y more than image width, height")
	}

	if width == 0 || imCols < uint(x)+width {
		width = imCols - uint(x)
	}

	if height == 0 || imRows < uint(y)+height {
		height = imRows - uint(y)
	}
	return mgw.CropImage(width, height, x, y)
}

func Composite(mgw *imagick.MagickWand, blob *[]byte, x, y int) error {
	source := imagick.NewMagickWand()
	source.ReadImageBlob(*blob)
	return mgw.CompositeImage(source, imagick.COMPOSITE_OP_SRC_IN, true, x, y)
}
