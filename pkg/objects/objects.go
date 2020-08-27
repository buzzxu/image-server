package objects

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
)

type (
	//字体
	Font struct {
		Size  float64 `json:"size"`
		Color string  `json:"color"`
		Align string  `json:"align"`
		Name  string  `json:"name"`
		DPI   float64 `json:"dpi"`
	}
)

func (f *Font) Gravity() imagick.GravityType {
	gravity := strings.ToLower(f.Align)
	switch gravity {
	case "center":
		return imagick.GRAVITY_CENTER
	case "left":
		return imagick.GRAVITY_WEST
	case "right":
		return imagick.GRAVITY_EAST
	case "top":
		return imagick.GRAVITY_NORTH
	case "bottom":
		return imagick.GRAVITY_SOUTH
	case "northwest":
		return imagick.GRAVITY_NORTH_WEST
	case "northeast":
		return imagick.GRAVITY_NORTH_EAST
	case "southwest":
		return imagick.GRAVITY_SOUTH_WEST
	case "southeast":
		return imagick.GRAVITY_SOUTH_EAST
	}
	return imagick.GRAVITY_CENTER
}
