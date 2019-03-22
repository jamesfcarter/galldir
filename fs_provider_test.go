package galldir_test

import (
	"io/ioutil"
	"testing"

	"github.com/jamesfcarter/galldir"
)

func TestReadDir(t *testing.T) {
	files, err := galldir.FsBackend("testdata").ReadDir("/")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatal("unexpected number of files")
	}
	for i, fn := range []string{"album", "hello"} {
		if files[i].Name() != fn {
			t.Errorf(`unexpected file "%s"`, files[i].Name())
		}
	}
}

func TestOpen(t *testing.T) {
	r, err := galldir.FsBackend("testdata").Open("hello")
	if err != nil {
		t.Fatal(err)
	}
	hello, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(hello) != "hello" {
		t.Errorf(`expected "hello", got "%s"`, string(hello))
	}
}

func TestPath(t *testing.T) {
	tests := []struct{ dir, path, expected string }{
		{"/pics", "fred", "/pics/fred"},
		{"/pics", "../fred", "/pics"},
		{"/pics", "/jim/../fred", "/pics/fred"},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			p := galldir.FsBackend(tc.dir).Path(tc.path)
			if p != tc.expected {
				t.Errorf(`expected "%s", got "%s"`, tc.expected, p)
			}
		})
	}
}
