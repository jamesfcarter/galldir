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
	"net/http"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	cache "github.com/patrickmn/go-cache"
)

// Provider is used to fetch Albums and Images from a Backend
type Provider struct {
	FS    http.FileSystem
	Cache *cache.Cache
}

func NewProvider(backend http.FileSystem) *Provider {
	return &Provider{
		FS:    backend,
		Cache: cache.New(0, 0),
	}
}

func (p *Provider) loadFile(path string) string {
	cacheName := CacheName("file", path)
	cacheVal, cached := p.Cache.Get(cacheName)
	if cached {
		return cacheVal.(string)
	}
	var content string
	f, err := p.FS.Open(path)
	if err == nil {
		contentBytes, _ := ioutil.ReadAll(f)
		content = string(contentBytes)
	}
	content = strings.TrimSuffix(content, "\n")
	p.Cache.Set(cacheName, content, cache.DefaultExpiration)
	return content
}

func (p *Provider) getName(path string) string {
	if name := p.loadFile(filepath.Join(path, ".title")); name != "" {
		return name
	}
	return NameFromPath(path)
}

// Album retrieves an Album from the backend, or returns an error if
// it is unable to.
func (p *Provider) Album(path string) (*Album, error) {
	cacheName := CacheName("album", path)
	cacheVal, cached := p.Cache.Get(cacheName)
	if cached {
		return cacheVal.(*Album), nil
	}
	album, err := p.loadAlbum(path)
	if err != nil {
		return nil, err
	}
	p.Cache.Set(cacheName, album, cache.DefaultExpiration)
	return album, nil
}

func (p *Provider) loadAlbum(path string) (*Album, error) {
	albumFile, err := p.FS.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := albumFile.Stat()
	if err != nil {
		return nil, err
	}
	a := &Album{
		Path: path,
		Name: p.getName(path),
		Time: fi.ModTime(),
	}
	files, err := albumFile.Readdir(0)
	if err != nil {
		return nil, err
	}
	a.Images = make([]Image, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && !IsImage(file.Name()) {
			continue
		}
		path := filepath.Join(path, file.Name())
		a.Images = append(a.Images, Image{
			Path: path,
			Name: func() string {
				if file.IsDir() {
					return p.getName(path)
				}
				return file.Name()
			}(),
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
	return p.FS.Open(path)
}

func (p *Provider) resizedImage(src io.ReadSeeker, width int, cacheName string) (io.ReadSeeker, error) {
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

// CacheName generates a unique key for the cache
func CacheName(class, path string) string {
	return class + "-" + path
}

// ThumbName generates a unique key for a thumbname in the cache
func ThumbName(class string, width int, path string) string {
	return CacheName(fmt.Sprintf("%s%d", class, width), path)
}

// CachedThumb returns a (potentially cached) thumbnail of the supplied
// source image
func (p *Provider) CachedThumb(cacheName string, width int, src io.ReadSeeker) (io.ReadSeeker, error) {
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	return p.resizedImage(src, width, cacheName)
}

// ImageThumb returns a (potentially cached) thumbnail of the image
// at the given path, scaled to the width.
func (p *Provider) ImageThumb(path string, width int) (io.ReadSeeker, error) {
	cacheName := ThumbName("thumb", width, path)
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	src, err := p.ImageContent(path)
	if err != nil {
		return nil, err
	}
	return p.resizedImage(src, width, cacheName)
}

// AlbumThumb  returns a (potentially cached) thumbnail of the album cover,
// scaled to the width
func (p *Provider) CoverThumb(album *Album, width int) (io.ReadSeeker, error) {
	if cover := p.loadFile(filepath.Join(album.Path, ".cover")); cover != "" {
		return p.ImageThumb(filepath.Join(album.Path, cover), width)
	}
	photos := album.Photos()
	if len(photos) == 0 {
		return nil, errors.New("no photo for cover " + album.Path)
	}
	for i := range photos {
		if strings.Contains(photos[i].Name, "star") {
			return p.ImageThumb(photos[i].Path, width)
		}
	}
	return p.ImageThumb(photos[0].Path, width)
}
