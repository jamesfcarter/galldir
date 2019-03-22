package galldir

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FsBackend satisfies the Backend interface and can be used to access
// files and and subdirectories of the supplied directory.
type FsBackend string

// ReadDir reads the directory at the supplied path and returns a list of
// directory entries sorted by filename.
func (backend FsBackend) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(backend.Path(path))
}

// Open opens the file at the supplied path for reading.
func (backend FsBackend) Open(path string) (io.ReadSeeker, error) {
	return os.Open(backend.Path(path))
}

// Path returns the full absolute path name of the supplied path. If
// an attempt is made to escape the backend root using .. then the
// root is returned.
func (backend FsBackend) Path(path string) string {
	root := string(backend)
	path = filepath.Clean(path)
	if path == ".." || strings.HasPrefix(path, "../") {
		return root
	}
	return filepath.Join(root, path)
}
