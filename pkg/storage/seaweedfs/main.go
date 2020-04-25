package seaweedfs

import (
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	ParamCollection   = "collection"
	ParamTTL          = "ttl"
	HeadAuthorization = "Authorization"
)

var seaweed *Seaweed

type Seaweedfs struct {
}

func (image *Seaweedfs) Init() {
	var token string
	if len(conf.Config.Seaweed.Secret) > 0 {
		claims := &jwt.StandardClaims{
			Issuer: "image-server",
		}
		if ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.Config.Seaweed.Secret)); err == nil {
			token = ss
		} else {
			log.Fatalf("jwt error %v", err)
		}
	}
	seaweed, _ = newSeaweed(conf.Config.Seaweed.MasterUrl, conf.Config.Seaweed.Filer, 8096, token, &http.Client{Timeout: 5 * time.Minute})
}

func (image *Seaweedfs) Check(params map[string]string) {
	if conf.Config.Seaweed.MasterUrl == "" {
		log.Printf("请设置MasterUrl")
	}
}
func (image *Seaweedfs) Upload(upload *storage.Upload) ([]string, error) {
	numfiles := len(upload.Blobs)
	paths := make([]string, numfiles)
	for index := 0; index < numfiles; index++ {
		fileName := upload.Keys[index]
		if upload.Rename {
			fileName = utils.NewFileName(upload.Folder, fileName)
		} else {
			fileName = filepath.Join(upload.Folder, fileName)
		}
		data := upload.Blobs[index]
		result, err := seaweed.filers[0].Upload(bytes.NewReader(*data), int64(len(*data)), fileName, "", "")
		if err != nil {
			return nil, err
		}

		if conf.Config.Domain != "" {
			paths[index] = conf.Config.Domain + fileName + "?id=" + result.FileID
		} else {
			paths[index] = fileName
		}

	}
	return paths, nil
}

func (image *Seaweedfs) Download(download *storage.Download) (*[]byte, string, error) {
	return nil, "", nil
}

func (image *Seaweedfs) Destory() {
	seaweed.close()
}
func (image *Seaweedfs) Delete(del *storage.Delete) (bool, error) {
	numfiles := len(del.Keys)
	for i := 0; i < numfiles; i++ {
		path := del.Keys[i]
		if err := seaweed.filers[0].Delete(path, nil); err != nil {
			return false, err
		}
	}
	return true, nil
}
