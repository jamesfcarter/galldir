package galldir

// Image specifies an image. An image may be the cover of a sub-album.
type Image struct {
	Path        string
	Name        string
	Description string
	IsAlbum     bool
}

// Album specifies a photo album
type Album struct {
	Path        string
	Name        string
	Description string
	Images      []Image
}
