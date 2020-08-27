package knife

import (
	"image-server/pkg/objects"
)

type (
	//裁剪
	CropParam struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	//合成
	CompositeParam struct {
		Text   string        `json:"text"`
		Url    string        `json:"url"`
		Base64 string        `json:"base_64"`
		Font   *objects.Font `json:"font"`
		X      int           `json:"x"`
		Y      int           `json:"y"`
		Width  int           `json:"width"`
		Height int           `json:"height"`
		Circle bool          `json:"circle"`
	}
)
