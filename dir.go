package galldir

import (
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

// Album specifies a photo album
type Album struct {
	Path        string
	Name        string
	Description string
	Images      []Image
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
