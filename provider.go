package galldir

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Backend interface {
	ReadDir(string) ([]os.FileInfo, error)
	Open(string) (io.Reader, error)
}

type Provider struct {
	Backend
}

func (p Provider) Album(path string) (*Album, error) {
	files, err := p.ReadDir(path)
	if err != nil {
		return nil, err
	}
	images := make([]Image, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && !isImage(file.Name()) {
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

func (p Provider) Image(path string) (io.Reader, error) {
	if !isImage(path) {
		return nil, errors.New("not an image")
	}
	return p.Open(path)
}

var imageExtensions = []string{".jpg", ".jpeg", ".png"}

func isImage(path string) bool {
	ext := filepath.Ext(path)
	for _, imExt := range imageExtensions {
		if strings.EqualFold(ext, imExt) {
			return true
		}
	}
	return false
}
