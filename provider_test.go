package galldir_test

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesfcarter/galldir"
)

func TestAlbum(t *testing.T) {
	tests := []struct {
		path      string
		name      string
		images    []string
		expectErr bool
	}{
		{
			path:   "/",
			name:   "Test Album",
			images: []string{"subalbum/"},
		},
		{
			path:      "/not_there",
			expectErr: true,
		},
		{
			path:   "/subalbum",
			name:   "Subalbum",
			images: []string{"icon.png"},
		},
	}
	provider := galldir.NewProvider(http.Dir("testdata/album"))
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
					if im.Path != filepath.Join(album.Path, filepath.Clean(name)) {
						t.Errorf(`unexpected image "%s"`, im.Name)
					}
					if im.IsAlbum != isAlbum {
						t.Error("incorrect isAlbum")
					}
				})
			}
			if tc.name != "" && album.Name != tc.name {
				t.Errorf("unexpected album name: %s", album.Name)
			}
		})
	}
}

func testHash(t *testing.T, r io.Reader, expected string) {
	if expected == "" {
		return
	}
	hasher := sha1.New()
	io.Copy(hasher, r)
	sha := fmt.Sprintf("%x", hasher.Sum(nil))
	if sha != expected {
		t.Errorf("bad hash: %s", sha)
	}
}

func TestImageContent(t *testing.T) {
	tests := []struct {
		path      string
		hash      string
		expectErr bool
	}{
		{
			path:      "/not_there.jpg",
			expectErr: true,
		},
		{
			path: "/subalbum/icon.png",
			hash: "aa72605dbcb4f8b933be68f0d11391673cd9ecc7",
		},
	}
	provider := galldir.NewProvider(http.Dir("testdata/album"))
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			r, err := provider.ImageContent(tc.path)
			if err == nil && tc.expectErr {
				t.Fatal("expected an error")
			}
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			testHash(t, r, tc.hash)
		})
	}
}

func TestImageThumb(t *testing.T) {
	tests := []struct {
		path      string
		hash      string
		expectErr bool
	}{
		{
			path:      "/not_there.jpg",
			expectErr: true,
		},
		{
			path: "/subalbum/icon.png",
			hash: "0a42d4b9ebda8bf872462e4c2a8c7934734f56b1",
		},
	}
	provider := galldir.NewProvider(http.Dir("testdata/album"))
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			r, err := provider.ImageThumb(tc.path, 10)
			if err == nil && tc.expectErr {
				t.Fatal("expected an error")
			}
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			testHash(t, r, tc.hash)
			// test again for cached copy
			r, err = provider.ImageThumb(tc.path, 10)
			if err != nil {
				t.Fatal(err)
			}
			testHash(t, r, tc.hash)
		})
	}
}

func TestCoverThumb(t *testing.T) {
	tests := []struct {
		path      string
		hash      string
		expectErr bool
	}{
		{
			path:      "/",
			expectErr: true,
		},
		{
			path: "/subalbum",
			hash: "0a42d4b9ebda8bf872462e4c2a8c7934734f56b1",
		},
	}
	provider := galldir.NewProvider(http.Dir("testdata/album"))
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			album, err := provider.Album(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			r, err := provider.CoverThumb(album, 10)
			if err == nil && tc.expectErr {
				t.Fatal("expected an error")
			}
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			testHash(t, r, tc.hash)
			// test again for cached copy
			r, err = provider.CoverThumb(album, 10)
			if err != nil {
				t.Fatal(err)
			}
			testHash(t, r, tc.hash)
		})
	}
}
