package galldir

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // loaded for image.Decode support
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	cache "github.com/patrickmn/go-cache"
)

const (
	fileTimeout  = 1 * time.Hour
	cachedImages = 64
)

// Provider is used to fetch Albums and Images from a Backend
type Provider struct {
	FS    http.FileSystem
	Cache *cache.Cache
	// ImageCacheEntries limits the number of cached full size images.
	// This is done to limit the amount of memory consumed, but also puts
	// an effective limit on the number of images that may be transferred
	// concurrently.
	ImageCacheEntries [cachedImages]struct {
		sync.Mutex
		Name string
	}
}

// NewProvider returns an initialized Provider
func NewProvider(backend http.FileSystem) *Provider {
	c := cache.New(0, 0)
	c.Set("imageIndex", uint(0), cache.NoExpiration)
	return &Provider{
		FS:    backend,
		Cache: c,
	}
}

func (p *Provider) claimCacheEntry(cacheName string, f func() error) error {
	index, err := p.Cache.IncrementUint("imageIndex", 1)
	if err != nil {
		panic(err)
	}
	index = index % cachedImages
	p.ImageCacheEntries[index].Lock()
	defer p.ImageCacheEntries[index].Unlock()
	if p.ImageCacheEntries[index].Name != "" {
		p.Cache.Delete(p.ImageCacheEntries[index].Name)
	}
	p.ImageCacheEntries[index].Name = cacheName
	return f()
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
	p.Cache.Set(cacheName, content, fileTimeout)
	return content
}

func (p *Provider) getDate(path string, modTime time.Time) time.Time {
	if date := p.loadFile(filepath.Join(path, ".date")); date != "" {
		t, err := time.Parse("2006-01-02 15:04:05", date)
		if err == nil {
			return t
		}
	}
	return modTime
}

func (p *Provider) getName(path string) string {
	if name := p.loadFile(filepath.Join(path, ".title")); name != "" {
		return name
	}
	return NameFromPath(path)
}

// Album retrieves a (possibly cached) Album from the backend, or returns an
// error if it is unable to.
func (p *Provider) Album(path string, refreshCache bool) (*Album, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	cacheName := CacheName("album", path)
	if !refreshCache {
		cacheVal, cached := p.Cache.Get(cacheName)
		if cached {
			return cacheVal.(*Album), nil
		}
	}
	album, err := p.loadAlbum(path)
	if err != nil {
		return nil, err
	}
	p.Cache.SetDefault(cacheName, album)
	return album, nil
}

func (p *Provider) loadAlbum(path string) (*Album, error) {
	albumFile, err := p.FS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %v", path, err)
	}
	fi, err := albumFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("Failed to stat %s: %v", path, err)
	}
	a := &Album{
		Path: path,
		Name: p.getName(path),
		Time: p.getDate(path, fi.ModTime()),
	}
	files, err := albumFile.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("Failed to readdir %s: %v", path, err)
	}
	a.Images = make([]Image, 0, len(files))
	for _, file := range files {
		fileName := strings.TrimPrefix(file.Name(), strings.TrimPrefix(path, "/"))
		if !file.IsDir() && !IsImage(fileName) {
			continue
		}
		path := filepath.Join(path, fileName)
		a.Images = append(a.Images, Image{
			Path: path,
			Name: func() string {
				if file.IsDir() {
					return p.getName(path)
				}
				return fileName
			}(),
			Time: func() time.Time {
				if file.IsDir() {
					return p.getDate(path, fi.ModTime())
				}
				return file.ModTime()
			}(),
			IsAlbum: file.IsDir(),
		})
	}
	return a, nil
}

// ImageContent returns an io.ReadSeeker for an image stored in the backend
// at the given path (that may have been cached). Any attempt to read
// anything other than an image will result in an error.
func (p *Provider) ImageContent(path string) (io.ReadSeeker, error) {
	if !IsImage(path) {
		return nil, errors.New("not an image")
	}
	cacheName := CacheName("image", path)
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	var image []byte
	err := p.claimCacheEntry(cacheName, func() error {
		src, err := p.FS.Open(path)
		if err != nil {
			return err
		}
		image, err = ioutil.ReadAll(src)
		if err != nil {
			return err
		}
		p.Cache.SetDefault(cacheName, image)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(image), nil
}

func (p *Provider) resizedImage(src io.ReadSeeker, size int, cacheName string) (io.ReadSeeker, error) {
	im, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	var thumb image.Image
	if p := im.Bounds().Size(); p.X > p.Y {
		thumb = imaging.Resize(im, size, 0, imaging.Lanczos)
	} else {
		thumb = imaging.Resize(im, 0, size, imaging.Lanczos)
	}
	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}
	jpgBytes := buf.Bytes()
	p.Cache.SetDefault(cacheName, jpgBytes)
	return bytes.NewReader(jpgBytes), nil
}

// CacheName generates a unique key for the cache
func CacheName(class, path string) string {
	return class + "-" + path
}

// ThumbName generates a unique key for a thumbname in the cache
func ThumbName(class string, size int, path string) string {
	return CacheName(fmt.Sprintf("%s%d", class, size), path)
}

// CachedThumb returns a (potentially cached) thumbnail of the supplied
// source image
func (p *Provider) CachedThumb(cacheName string, size int, src io.ReadSeeker) (io.ReadSeeker, error) {
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	return p.resizedImage(src, size, cacheName)
}

// ImageThumb returns a (potentially cached) thumbnail of the image
// at the given path, scaled to the size.
func (p *Provider) ImageThumb(path string, size int) (io.ReadSeeker, error) {
	cacheName := ThumbName("thumb", size, path)
	cachedImage, cached := p.Cache.Get(cacheName)
	if cached {
		return bytes.NewReader(cachedImage.([]byte)), nil
	}
	src, err := p.ImageContent(path)
	if err != nil {
		return nil, err
	}
	return p.resizedImage(src, size, cacheName)
}

// CoverThumb returns a (potentially cached) thumbnail of the album cover,
// scaled to the size
func (p *Provider) CoverThumb(album *Album, size int) (io.ReadSeeker, error) {
	if cover := p.loadFile(filepath.Join(album.Path, ".cover")); cover != "" {
		return p.ImageThumb(filepath.Join(album.Path, cover), size)
	}
	photos := album.Photos()
	if len(photos) == 0 {
		return nil, errors.New("no photo for cover " + album.Path)
	}
	return p.ImageThumb(photos[0].Path, size)
}
