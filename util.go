package galldir

import (
	"path/filepath"
	"regexp"
	"strings"
)

var imageExtensions = []string{".jpg", ".jpeg", ".png"}

// IsImage takes a path name and returns true if it refers to an image file
func IsImage(path string) bool {
	ext := filepath.Ext(path)
	for _, imExt := range imageExtensions {
		if strings.EqualFold(ext, imExt) {
			return true
		}
	}
	return false
}

var nameRegexp = []struct {
	r *regexp.Regexp
	s string
}{
	{regexp.MustCompile(`(\d\d\d\d)_(\d\d)_(\d\d)`), "$1-$2-$3"},
	{regexp.MustCompile(`_`), " "},
}

// NameFromPath tries to generate a human friendly name from a path
func NameFromPath(path string) string {
	str := []byte(filepath.Base(path))
	for _, re := range nameRegexp {
		str = re.r.ReplaceAll(str, []byte(re.s))
	}
	return strings.Title(string(str))
}
