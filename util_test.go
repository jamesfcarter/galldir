package galldir_test

import (
	"testing"

	"github.com/jamesfcarter/galldir"
)

func TestIsImage(t *testing.T) {
	tests := []struct {
		path    string
		isImage bool
	}{
		{"/foo", false},
		{"/foo.jpg/x", false},
		{"/foo/x.jpg", true},
		{"/foo/x.JPG", true},
		{"/foo/x.png", true},
		{"/foo/x.PnG", true},
		{"/foo/.jpeg", true},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			if galldir.IsImage(tc.path) != tc.isImage {
				t.Fatal("failed")
			}
		})
	}
}
