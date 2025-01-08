package knife

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	"image-server/pkg/imagemagick"
	"image-server/pkg/img"
	"image-server/pkg/utils"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func init() {
	imagick.Initialize()
}

// 裁剪
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
		_blob, err := mw.GetImageBlob()
		if err != nil {
			return nil, err
		} else {
			data[i] = prefix + base64.StdEncoding.EncodeToString(_blob)
		}

		mw.Clear()
	}
	return data, nil
}

// 合成
func Composite0(width, height uint, composite *CompositeParam) (*[]byte, error) {
	mw := imagick.NewMagickWand()
	pw := imagick.NewPixelWand()
	mw.NewImage(width, height, pw)
	defer mw.Destroy()
	defer pw.Destroy()
	mw.StripImage()
	if composite.Text != "" {
		//增加文字
		imagemagick.AppendText(mw, float64(composite.X), float64(composite.Y), composite.Text, composite.Font)
	} else {
		var blob *[]byte
		var err error
		if composite.Url != "" {
			blob, err = utils.GetUrlBuffer(composite.Url)
			if err != nil {
				return nil, err
			}
		} else if composite.Base64 != "" {
			image := composite.Base64
			index := strings.Index(composite.Base64, ",")
			if index > 0 {
				image = composite.Base64[index+1:]
			}
			bytes, err := base64.StdEncoding.DecodeString(image)
			if err != nil {
				return nil, err
			}
			blob = &bytes
		} else {
			return nil, nil
		}
		imagemagick.Composite(mw, blob, composite.Width, composite.Width)
	}

	return nil, nil
}

// 合成
func Composite(composites []CompositeParam) (*[]byte, error) {
	var (
		w = 0
		h = 0
	)
	for i := 0; i < len(composites); i++ {
		if w < composites[i].Width {
			w = composites[i].Width
		}
		if h < composites[i].Height {
			h = composites[i].Height
		}
	}
	des := image.NewRGBA(image.Rect(0, 0, w, h))
	for _, composite := range composites {
		if composite.Text != "" {
			if err := img.MergeText(des, composite.Text, composite.X, composite.Y, composite.Font); err != nil {
				return nil, err
			}
		} else {
			var blob *[]byte
			var err error
			if composite.Url != "" {
				blob, err = utils.GetUrlBuffer(composite.Url)
				if err != nil {
					return nil, err
				}
			} else if composite.Base64 != "" {
				image := composite.Base64
				index := strings.Index(composite.Base64, ",")
				if index > 0 {
					image = composite.Base64[index+1:]
				}
				bytes, err := base64.StdEncoding.DecodeString(image)
				if err != nil {
					return nil, err
				}
				blob = &bytes
			} else {
				return nil, nil
			}
			target, err := utils.Byte2Image(blob)
			if err != nil {
				return nil, err
			}
			if w != composite.Width && h != composite.Height {
				target = imaging.Resize(target, composite.Width, composite.Height, imaging.Lanczos)
			}

			//result = imaging.Paste(des,img,image.Pt(int(composite.X),int(composite.Y)))
			if composite.Circle {
				img.DrawCircle(des, target, composite.X, composite.Y)
			} else {
				draw.Draw(des, des.Bounds().Add(image.Pt(composite.X, composite.Y)), target, target.Bounds().Min, draw.Over)

			}
		}
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, des); err != nil {
		return nil, err
	}
	blob := buf.Bytes()
	return &blob, nil
}

func Convert(blobs []*[]byte, fileNames []string, target string, is_rename bool) (*[]byte, error) {
	return nil, nil
}
