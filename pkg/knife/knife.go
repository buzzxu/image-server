package knife

import (
	"encoding/base64"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image-server/pkg/imagemagick"
	"image-server/pkg/utils"
	_ "image/jpeg"
)

func init() {
	imagick.Initialize()
}

func Crop(blob *[]byte, crops *[]CropParam) ([]string, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	numParam := len(*crops)
	data := make([]string, numParam)
	prefix, err := utils.Base64ImagePrefix(blob)
	if err != nil {
		return nil, err
	}
	for i := 0; i < numParam; i++ {
		mw.ReadImageBlob(*blob)
		crop := (*crops)[i]
		error := imagemagick.Crop(uint(crop.Width), uint(crop.Height), crop.X, crop.Y, mw)
		if error != nil {
			return nil, error
		}
		data[i] = prefix + base64.StdEncoding.EncodeToString(mw.GetImageBlob())
		mw.Clear()
	}
	return data, nil
}
