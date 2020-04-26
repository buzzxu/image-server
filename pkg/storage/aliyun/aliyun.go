package aliyun

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/buzzxu/boys/types"
	"image-server/pkg/conf"
	"image-server/pkg/storage"
	"image-server/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	client     *oss.Client
	bucket     *oss.Bucket
	url        string
	defaultImg *[]byte
	urlLength  int
)

type Aliyun struct {
}

func (image *Aliyun) Init() {
	var err error
	if client, err = oss.New(conf.Config.Aliyun.Endpoint, conf.Config.Aliyun.AccessKeyId, conf.Config.Aliyun.AccessKeySecret); err != nil {
		log.Fatalf("阿里云")
	}
	if bucket, err = client.Bucket(conf.Config.Aliyun.Bucket); err != nil {
		log.Printf("获取Bucket[%s] 失败，需要重建", conf.Config.Aliyun.Bucket)
		if err = client.CreateBucket(conf.Config.Aliyun.Bucket); err != nil {
			log.Fatalf("创建Bucket[%s] 失败,原因 %s", conf.Config.Aliyun.Bucket, err.Error())
		}
		if bucket, err = client.Bucket(conf.Config.Aliyun.Bucket); err != nil {
			log.Fatalf("重新获取Bucket[%s] 失败,原因 %s", conf.Config.Aliyun.Bucket, err.Error())
		}
	}
	if conf.Config.Domain == "" {
		url = conf.Config.Aliyun.Endpoint[:8] + conf.Config.Aliyun.Bucket + "." + conf.Config.Aliyun.Endpoint[8:]
	} else {
		url = conf.Config.Domain
	}
	if defaultImg, err = storage.GetDefaultImg(); err != nil {
		log.Fatalf("默认图片加载失败,原因 %s", err.Error())
	}
	urlLength = len(url) + 1
}

func (image *Aliyun) Check(params map[string]string) {
	if bucket == nil {
		log.Fatalf("无法获取Bucket[%s]", conf.Config.Aliyun.Bucket)
	}
}

func (image *Aliyun) Upload(upload *storage.Upload) ([]string, error) {
	numfiles := len(upload.Blobs)
	paths := make([]string, numfiles)
	for index := 0; index < numfiles; index++ {
		fileName := upload.Keys[index]
		if upload.Rename {
			fileName = utils.NewFileName(upload.Folder, fileName)
		} else {
			fileName = filepath.Join(upload.Folder, fileName)
		}
		if err := bucket.PutObject(strings.TrimPrefix(fileName, "/"), bytes.NewReader(*upload.Blobs[index])); err != nil {
			return nil, err
		}
		paths[index] = url + fileName
	}
	return paths, nil
}

func (image *Aliyun) Download(download *storage.Download) (*[]byte, string, error) {
	blob, err := get(download.Path[1:])
	return blob, http.DetectContentType(*blob), err
}

func (image *Aliyun) Delete(del *storage.Delete) (bool, error) {
	var err error
	numfiles := len(del.Keys)
	if numfiles == 1 {
		err = bucket.DeleteObject(del.Keys[0][urlLength:])
	} else {
		var keys = make([]string, len(del.Keys))
		for index := 0; index < numfiles; index++ {
			keys[index] = del.Keys[index][urlLength:]
		}
		_, err = bucket.DeleteObjects(keys)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (image *Aliyun) Destory() {

}

//从阿里云获取数据
func get(key string) (*[]byte, error) {
	read, err := bucket.GetObject(key)
	if err != nil {
		return defaultImg, types.ErrNotFound
	}
	blob, err := ioutil.ReadAll(read)
	if err != nil {
		return defaultImg, types.ErrorOf(err)
	}
	return &blob, nil
}

/*func uploadDirToS3(dir string, svc *s3.S3) {
	fileList := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	var wg sync.WaitGroup
	wg.Add(len(fileList))
	for _, pathOfFile := range fileList[1:] {
		//maybe spin off a goroutine here??
		go putInS3(pathOfFile, svc, &wg)
	}
	wg.Wait()
}

func putInS3(pathOfFile string, svc *s3.S3, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	file, _ := os.Open(pathOfFile)
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := file.Name()
	params := &s3.PutObjectInput{
		Bucket:        aws.String("bucket-name"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, _ := svc.PutObject(params)
	fmt.Printf("response %s", awsutil.StringValue(resp))
}*/
