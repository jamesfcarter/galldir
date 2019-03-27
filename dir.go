package galldir

import (
	"sort"
	"strings"
	"time"
)

// Image specifies an image. An image may be the cover of a sub-album.
type Image struct {
	Path        string
	Name        string
	Description string
	Time        time.Time
	IsAlbum     bool
}

type ImagesByName []Image

func (a ImagesByName) Len() int      { return len(a) }
func (a ImagesByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ImagesByName) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
}

type ImagesByTime []Image

func (a ImagesByTime) Len() int           { return len(a) }
func (a ImagesByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ImagesByTime) Less(i, j int) bool { return a[j].Time.Before(a[i].Time) }

// Album specifies a photo album
type Album struct {
	Path        string
	Name        string
	Description string
	Images      []Image
	Time        time.Time
}

// Image returns an Image from an Album or nil if it cannot be found
func (a *Album) Image(path string) *Image {
	for i := range a.Images {
		if a.Images[i].Path == path {
			return &a.Images[i]
		}
	}
	return nil
}

func (a *Album) images(isAlbum bool) []Image {
	result := make([]Image, 0, len(a.Images))
	for _, im := range a.Images {
		if im.IsAlbum != isAlbum {
			continue
		}
		result = append(result, im)
	}
	return result
}

// Images returns a list of images from an album that are not sub-albums
func (a *Album) Photos() []Image {
	images := a.images(false)
	sort.Sort(ImagesByName(images))
	return images
}

// Albums returns a list of images from an album that are sub-albums
func (a *Album) Albums() []Image {
	images := a.images(true)
	sort.Sort(ImagesByTime(images))
	return images
}
