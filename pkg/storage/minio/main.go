package minio

import (
	"bytes"
	"github.com/buzzxu/boys/types"
	"github.com/minio/minio-go/v6"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Minio struct {
	minioClient  *minio.Client
	bucket       string
	cacheControl string
	defaultImg   *[]byte
}

func (image *Minio) Init() {
	var err error
	image.minioClient, err = minio.New(conf.Config.Minio.Endpoint, conf.Config.Minio.AccessKey, conf.Config.Minio.SecretKey, conf.Config.Minio.UseSSL)
	if err != nil {
		log.Fatalf("minio initialize error %v", err)
	}
	image.cacheControl = "public,max-age=" + strconv.Itoa(conf.Config.MaxAge)
	image.bucket = conf.Config.Minio.Bucket
	if image.defaultImg, err = storage.GetDefaultImg(); err != nil {
		log.Fatalf("默认图片加载失败,原因 %s", err.Error())
	}
}

func (image *Minio) Check(params map[string]string) {
	exists, err := image.minioClient.BucketExists(conf.Config.Minio.Bucket)
	if err == nil && exists {
		policy, err := image.minioClient.GetBucketPolicy(conf.Config.Minio.Bucket)
		if err != nil {
			log.Fatalf("get bucket policy err %v", err)
		}
		log.Printf("found bucket %s ", conf.Config.Minio.Bucket)
		log.Printf("Policy %s \n", policy)
	} else {
		if err := image.minioClient.MakeBucket(conf.Config.Minio.Bucket, conf.Config.Minio.Location); err != nil {
			log.Fatalf("make bucket err %v", err)
		}
		var policy = strings.ReplaceAll("{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:ListBucket\",\"s3:GetBucketLocation\"],\"Resource\":[\"arn:aws:s3:::{bucket}\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::{bucket}/*\"]}]}", "{bucket}", conf.Config.Minio.Bucket)
		log.Printf("Set Policy:%s\n", policy)
		//read ,write,read_write
		image.minioClient.SetBucketPolicy(conf.Config.Minio.Bucket, policy)
		log.Printf("Successfully created %s", conf.Config.Minio.Bucket)
	}
}
func (image *Minio) Upload(upload *storage.Upload) ([]string, error) {
	//TODO 支持Webp
	numfiles := len(upload.Blobs)
	paths := make([]string, numfiles)
	for index := 0; index < numfiles; index++ {
		fileName := upload.Keys[index]
		if upload.Rename {
			fileName = utils.NewFileName(upload.Folder, fileName)
		} else {
			fileName = filepath.Join(upload.Folder, fileName)
		}
		data := *upload.Blobs[index]

		_, err := image.minioClient.PutObject(image.bucket, strings.TrimPrefix(fileName, "/"), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
			CacheControl: image.cacheControl,
			ContentType:  http.DetectContentType(data),
		})
		if err != nil {
			return nil, err
		}
		if conf.Config.Domain != "" {
			paths[index] = conf.Config.Domain + fileName
		} else {
			paths[index] = fileName
		}

	}
	return paths, nil
}
func (image *Minio) Delete(del *storage.Delete) (bool, error) {
	numfiles := len(del.Keys)
	if numfiles == 1 {
		if err := image.minioClient.RemoveObject(image.bucket, strings.TrimPrefix(del.Keys[0], "/")); err != nil {
			return false, err
		}
	} else {
		ch := make(chan string, numfiles)

		go func() {
			defer close(ch)
			for i := 0; i < numfiles; i++ {
				ch <- strings.TrimPrefix(del.Keys[i], "/")
			}
		}()
		for err1 := range image.minioClient.RemoveObjects(image.bucket, ch) {
			return false, err1.Err
		}
	}
	return true, nil
}

func (image *Minio) Download(download *storage.Download) (*[]byte, string, error) {
	object, err := image.minioClient.GetObject(image.bucket, download.FileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}
	blob, err := ioutil.ReadAll(object)
	if err != nil {
		return image.defaultImg, conf.Config.DefaultImg, types.ErrorOf(err)
	}
	return &blob, http.DetectContentType(blob), nil
}

func (image *Minio) Destory() {
}
