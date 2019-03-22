package galldir

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FsBackend string

func (backend FsBackend) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(backend.Path(path))
}

func (backend FsBackend) Open(path string) (io.Reader, error) {
	return os.Open(backend.Path(path))
}

func (backend FsBackend) Path(path string) string {
	root := string(backend)
	path = filepath.Clean(path)
	if path == ".." || strings.HasPrefix(path, "../") {
		return root
	}
	return filepath.Join(root, path)
}
