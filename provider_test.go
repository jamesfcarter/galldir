package galldir_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesfcarter/galldir"
)

func TestAlbum(t *testing.T) {
	tests := []struct {
		path      string
		images    []string
		expectErr bool
	}{
		{"/", []string{"subalbum/"}, false},
		{"/not_there", []string{}, true},
		{"/subalbum", []string{"icon.png"}, false},
	}
	provider := galldir.Provider{galldir.FsBackend("testdata/album")}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			album, err := provider.Album(tc.path)
			if err == nil && tc.expectErr {
				t.Fatal("expected an error")
			}
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			imageCount := len(album.Images)
			if imageCount != len(tc.images) {
				t.Fatalf("unexpected number of images: %d", imageCount)
			}
			for i, name := range tc.images {
				t.Run(name, func(t *testing.T) {
					im := album.Images[i]
					isAlbum := strings.HasSuffix(name, "/")
					if im.Name != filepath.Clean(name) {
						t.Errorf(`unexpected image "%s"`, im.Name)
					}
					if im.IsAlbum != isAlbum {
						t.Error("incorrect isAlbum")
					}
				})
			}
		})
	}
}
