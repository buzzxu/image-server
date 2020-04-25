package seaweedfs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Filer struct {
	base   *url.URL
	client *httpClient
	token  string
}

type FilerUploadResult struct {
	Name    string `json:"name,omitempty"`
	FileURL string `json:"url,omitempty"`
	FileID  string `json:"fid,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewFiler(u, token string, client *http.Client) (f *Filer, err error) {
	return newFiler(u, token, newHTTPClient(client))
}

func newFiler(u, token string, client *httpClient) (f *Filer, err error) {
	base, err := parseURI(u)
	if err != nil {
		return
	}

	f = &Filer{
		base:   base,
		client: client,
		token:  token,
	}

	return
}
func (f *Filer) Close() (err error) {
	if f.client != nil {
		err = f.client.Close()
	}
	return
}

// UploadFile a file.
func (f *Filer) UploadFile(localFilePath, newPath, collection, ttl string) (result *FilerUploadResult, err error) {
	fp, err := NewFilePart(localFilePath)
	if err == nil {
		var data []byte
		data, _, err = f.client.upload(encodeURI(*f.base, newPath, normalize(nil, collection, ttl)), localFilePath, fp.Reader, fp.MimeType, f.token)
		if err == nil {
			result = &FilerUploadResult{}
			err = json.Unmarshal(data, result)
		}

		_ = fp.Close()
	}
	return
}

// Upload content.
func (f *Filer) Upload(content io.Reader, fileSize int64, newPath, collection, ttl string) (result *FilerUploadResult, err error) {
	fp := NewFilePartFromReader(ioutil.NopCloser(content), newPath, fileSize)

	var data []byte
	data, _, err = f.client.upload(encodeURI(*f.base, newPath, normalize(nil, collection, ttl)), newPath, ioutil.NopCloser(content), "", f.token)
	if err == nil {
		result = &FilerUploadResult{}
		err = json.Unmarshal(data, result)
	}

	_ = fp.Close()

	return
}

// Get response data from filer.
func (f *Filer) Get(path string, args url.Values, header map[string]string) (data []byte, statusCode int, err error) {
	data, statusCode, err = f.client.get(encodeURI(*f.base, path, args), header)
	return
}

// Download a file.
func (f *Filer) Download(path string, args url.Values, callback func(io.Reader) error) (err error) {
	_, err = f.client.download(encodeURI(*f.base, path, args), callback)
	return
}

// Delete a file/dir.
func (f *Filer) Delete(path string, args url.Values) (err error) {
	_, err = f.client.delete(encodeURI(*f.base, path, args), f.token)
	return
}
