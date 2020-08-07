package none

import "image-server/pkg/storage"

type None struct {
}

func (image *None) Init() {

}

func (image *None) Check(params map[string]string) {

}

func (image *None) Upload(upload *storage.Upload) ([]string, error) {
	return nil, nil
}

func (image *None) Download(download *storage.Download) (*[]byte, string, error) {
	return nil, "", nil
}

func (image *None) Delete(del *storage.Delete) (bool, error) {
	return false, nil
}

func (image *None) Destory() {

}
