package galldir

type Image struct {
	Path        string
	Name        string
	Description string
	IsAlbum     bool
}

type Album struct {
	Path        string
	Name        string
	Description string
	Images      []Image
}
