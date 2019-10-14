package aliyun

import "image-server/pkg/storage"

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
	return nil, nil
}

func (image *Aliyun) Destory() {

}
