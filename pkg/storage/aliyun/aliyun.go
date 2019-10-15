package aliyun

import (
	"image-server/pkg/storage"
)

type Aliyun struct {
}

func (image *Aliyun) Init() {

}

func (image *Aliyun) Check(params map[string]string) {

}

func (image *Aliyun) Upload(upload *storage.Upload) ([]string, error) {
	return nil, nil
}

func (image *Aliyun) Download(download *storage.Download) ([]byte, string, error) {
	return nil, "", nil
}

func (image *Aliyun) Destory() {

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
