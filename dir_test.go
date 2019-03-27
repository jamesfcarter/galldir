package galldir_test

import (
	"testing"
	"time"

	"github.com/jamesfcarter/galldir"
)

var testAlbum = &galldir.Album{
	Path: "/",
	Images: []galldir.Image{
		galldir.Image{
			Path:    "/foo",
			Name:    "foo",
			Time:    time.Unix(1234, 0),
			IsAlbum: true,
		},
		galldir.Image{
			Path: "/bar.jpg",
			Name: "bar",
		},
		galldir.Image{
			Path: "/baz.png",
			Name: "baz",
		},
		galldir.Image{
			Path:    "/qux",
			Name:    "qux",
			Time:    time.Unix(5678, 0),
			IsAlbum: true,
		},
	},
}

func TestImage(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"/baz.png", "baz"},
		{"/foo", "foo"},
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

func testImages(t *testing.T, result []galldir.Image, name []string) {
	if count := len(result); count != len(name) {
		t.Fatalf("unexpected number of results: %d", count)
	}
	for i := range result {
		if rname := result[i].Name; rname != name[i] {
			t.Errorf("unexpected result: %s", rname)
		}
	}
}

func TestPhotos(t *testing.T) {
	testImages(t, testAlbum.Photos(), []string{"bar", "baz"})
}

func TestAlbums(t *testing.T) {
	testImages(t, testAlbum.Albums(), []string{"qux", "foo"})
}
