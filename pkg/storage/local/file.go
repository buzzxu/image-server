package local

import (
	"errors"
	"image-server/pkg/storage"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	root string
	key  string
}

func (f *File) Key() string {
	return f.key
}

func (f *File) Exist() (bool, string, error) {
	_, err := os.Stat(filepath.Join(f.root, f.key))
	if err == nil {
		return true, "", nil
	}
	if os.IsNotExist(err) {
		return false, "", nil
	}
	return false, "", err
}

func (f *File) Meta() (storage.Fetcher, error) {
	return nil, errors.New("Meta not supported")
}

func (f *File) Append(blob []byte, index int64, kvs ...storage.KV) (int64, string, error) {
	return index, "", errors.New("Append not supported")
}

func (f *File) Delete() (string, error) {
	return "", os.Remove(f.Key())
}

func (f *File) Bytes() ([]byte, error) {
	blob, err := ioutil.ReadFile(filepath.Join(f.root, f.key))
	return blob, err
}

func (f *File) SetMeta(...storage.KV) error {
	return errors.New("SetMeta not supported")
}
