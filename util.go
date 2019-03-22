package galldir

import (
	"path/filepath"
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
