package galldir_test

import (
	"testing"

	"github.com/jamesfcarter/galldir"
)

var foo = galldir.Image{
	Path: "/foo.png",
	Name: "foo",
}

var bar = galldir.Image{
	Path:    "/bar",
	Name:    "bar",
	IsAlbum: true,
}

var testAlbum = &galldir.Album{
	Path:   "/",
	Images: []galldir.Image{foo, bar},
}

func TestImage(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"/foo.png", "foo"},
		{"/bar", "bar"},
		{"/baz", ""},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			var name string
			im := testAlbum.Image(tc.path)
			if im != nil {
				name = im.Name
			}
			if tc.result != name {
				t.Errorf("unexpected image found: %s", name)
			}
		})
	}
}

func testImages(t *testing.T, result []galldir.Image, name string) {
	if count := len(result); count != 1 {
		t.Fatalf("unexpected number of results: %d", count)
	}
	if rname := result[0].Name; rname != name {
		t.Errorf("unexpected result: %s", rname)
	}
}

func TestPhotos(t *testing.T) {
	testImages(t, testAlbum.Photos(), "foo")
}

func TestAlbums(t *testing.T) {
	testImages(t, testAlbum.Albums(), "bar")
}
