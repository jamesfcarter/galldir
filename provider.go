package galldir

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	cache "github.com/patrickmn/go-cache"
)

// Backend is an interface for accessing images and the folders in
// which they are stored. Access is read only.
type Backend interface {
	ReadDir(string) ([]os.FileInfo, error)
	Open(string) (io.ReadSeeker, error)
}

// Provider is used to fetch Albums and Images from a Backend
type Provider struct {
	Backend
	Cache *cache.Cache
}

func NewProvider(b Backend) *Provider {
	return &Provider{
		Backend: b,
		Cache:   cache.New(0, 0),
	}
}

func (p *Provider) loadFile(a *Album, name string) string {
	path := filepath.Join(a.Path, name)
	f, err := p.Open(path)
	if err != nil {
		return ""
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(content)
}

func (p *Provider) getAlbumName(a *Album) {
	name := p.loadFile(a, ".title")
	if name == "" {
		name = filepath.Base(a.Path)
	}
	a.Name = name
}

// Album retrieves an Album from the backend, or returns an error if
// it is unable to.
func (p *Provider) Album(path string) (*Album, error) {
	a := &Album{Path: path}
	p.getAlbumName(a)
	files, err := p.ReadDir(path)
	if err != nil {
		return nil, err
	}
	a.Images = make([]Image, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && !IsImage(file.Name()) {
			continue
		}
		a.Images = append(a.Images, Image{
			Path:    filepath.Join(path, file.Name()),
			Name:    file.Name(),
			Time:    file.ModTime(),
			IsAlbum: file.IsDir(),
		})
	}
	return a, nil
}

// ImageContent returns an io.ReadSeeker for an image stored in the backend
// at the given path. Any attempt to read anything other than an image will
// result in an error.
func (p *Provider) ImageContent(path string) (io.ReadSeeker, error) {
	if !IsImage(path) {
		return nil, errors.New("not an image")
	}
	return p.Open(path)
}

// ImageThumb returns a (potentially cached) thumbnail of the image
// at the given path, scaled to the width.
func (p *Provider) ImageThumb(path string, width int) (io.ReadSeeker, error) {
	cacheName := fmt.Sprintf("thumb%d-%s", width, path)
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	src, err := p.ImageContent(path)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}
	thumb := imaging.Resize(im, width, 0, imaging.Lanczos)
	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, err
	}
	jpgBytes := buf.Bytes()
	p.Cache.Set(cacheName, jpgBytes, cache.DefaultExpiration)
	return bytes.NewReader(jpgBytes), nil
}
