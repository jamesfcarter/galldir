package galldir

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Backend is an interface for accessing images and the folders in
// which they are stored. Access is read only.
type Backend interface {
	ReadDir(string) ([]os.FileInfo, error)
	Open(string) (io.Reader, error)
}

// Provider is used to fetch Albums and Images from a Backend
type Provider struct {
	Backend
}

// Album retreives an Album from the backend, or returns an error if
// it is unable to.
func (p Provider) Album(path string) (*Album, error) {
	files, err := p.ReadDir(path)
	if err != nil {
		return nil, err
	}
	images := make([]Image, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && !IsImage(file.Name()) {
			continue
		}
		images = append(images, Image{
			Path:    filepath.Join(path, file.Name()),
			Name:    file.Name(),
			IsAlbum: file.IsDir(),
		})
	}
	return &Album{
		Path:   path,
		Images: images,
	}, nil
}

// Image returns an io.Reader for an image stored in the backend at the
// given path. Any attempt to read anything other than an image will result
// in an error.
func (p Provider) Image(path string) (io.Reader, error) {
	if !IsImage(path) {
		return nil, errors.New("not an image")
	}
	return p.Open(path)
}
