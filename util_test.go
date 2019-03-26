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

func TestNameFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"2018_01_01", "2018-01-01"},
		{"foo_bar_baz", "Foo Bar Baz"},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			if r := galldir.NameFromPath(tc.path); r != tc.expected {
				t.Fatalf("unexpected result: %s", r)
			}
		})
	}
}
