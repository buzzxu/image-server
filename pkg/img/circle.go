package img

import (
	"image"
	"image/color"
	"math"
)

type CircleMask struct {
	image    image.Image
	point    image.Point
	diameter int
}

func NewCircleMask(img image.Image, p image.Point, d int) CircleMask {
	return CircleMask{img, p, d}
}
func (ci CircleMask) ColorModel() color.Model {
	return ci.image.ColorModel()
}

func (ci CircleMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, ci.diameter, ci.diameter)
}
func (ci CircleMask) At(x, y int) color.Color {
	d := ci.diameter
	dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
	if dis > float64(d)/2 {
		return ci.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
	} else {
		return ci.image.At(ci.point.X+x, ci.point.Y+y)
	}
}

//func (ci CircleMask) At(target *image.RGBA,x, y ,diameter int) color.Color  {
//	dis := math.Sqrt(math.Pow(float64(x-diameter/2), 2) + math.Pow(float64(y-diameter/2), 2))
//	if dis > float64(diameter)/2 {
//		return target.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
//	} else {
//		return target.At(ci.point.X+x, ci.point.Y+y)
//	}
//}
