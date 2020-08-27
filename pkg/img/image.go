package img

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/labstack/gommon/log"
	"golang.org/x/image/font"
	"image"
	"image-server/pkg/objects"
	"image-server/pkg/utils"
	"image/draw"
	"io/ioutil"
	"os"
)

var fontType *truetype.Font

func init() {
	var err error
	fontType, err = LoadTextType("../../assets/msyh.ttf")
	if err != nil {
		currentDir, _ := os.Getwd()
		fontType, err = LoadTextType(currentDir + "/assets/msyh.ttf")
		if err != nil {
			fontType, err = LoadTextType("/app/msyh.ttf")
			if err != nil {
				log.Fatalf("读取配置文件内容失败,%v ", err)
			}
		}
	}
}

//新建图片载体
func NewPNG(X0 int, Y0 int, X1 int, Y1 int) *image.RGBA {
	return image.NewRGBA(image.Rect(X0, Y0, X1, Y1))
}

//合并图片到载体
func MergeImage(PNG draw.Image, image image.Image, imageBound image.Point) {
	draw.Draw(PNG, PNG.Bounds(), image, imageBound, draw.Over)
}

//圆形图片
func DrawCircle(backgroud draw.Image, target image.Image, x, y int) {
	// 算出图片的宽度和高试
	width := target.Bounds().Max.X - target.Bounds().Min.X
	hight := target.Bounds().Max.Y - target.Bounds().Min.Y
	srcPng := NewPNG(0, 0, width, hight)
	MergeImage(srcPng, target, target.Bounds().Min)
	diameter := width
	if width > hight {
		diameter = hight
	}
	// 遮罩
	srcMask := NewCircleMask(srcPng, image.Point{0, 0}, diameter)
	srcPoint := image.Point{
		X: x,
		Y: y,
	}
	MergeImage(backgroud, srcMask, target.Bounds().Min.Sub(srcPoint))
}

func MergeText(target *image.RGBA, text string, x, y int, f *objects.Font) error {
	c := freetype.NewContext()
	//设置屏幕每英寸的分辨率
	var dpi = f.DPI
	if dpi == 0 {
		dpi = 72
	}
	c.SetDPI(dpi)
	//设置剪裁矩形以进行绘制
	c.SetClip(target.Bounds())
	c.SetDst(target)
	c.SetHinting(font.HintingFull)
	color, err := utils.ParseHexColor(f.Color)
	if err != nil {
		return err
	}
	//设置绘制操作的源图像，通常为 image.Uniform
	c.SetSrc(image.NewUniform(color))
	//以磅为单位设置字体大小
	c.SetFontSize(f.Size)
	//设置用于绘制文本的字体
	c.SetFont(fontType)
	pt := freetype.Pt(x, y)
	_, err = c.DrawString(text, pt)
	return nil
}

func LoadTextType(path string) (*truetype.Font, error) {
	// 这里需要读取中文字体，否则中文文字会变成方格
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return &truetype.Font{}, err
	}

	return f, err
}
