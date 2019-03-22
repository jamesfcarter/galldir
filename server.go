package galldir

import (
"log"
	"net/http"
	"path"
	"fmt"
)

type Server struct {
	Provider Provider
}

func (s *Server) album(w http.ResponseWriter, r *http.Request) {
	album, err := s.Provider.Album(r.URL.Path)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "<html><body><ul>")
	for _, image := range album.Images {
	    fmt.Fprintf(w, "<li><a href='%s'>%s</a></li>", image.Path, image.Name)
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func (s *Server) image(w http.ResponseWriter, r *http.Request) {
	albumPath := path.Dir(r.URL.Path)
	album, err := s.Provider.Album(albumPath)
	if err != nil {
		log.Println(err)
		return
	}
	image := album.Image(r.URL.Path)
	if image == nil {
		log.Printf("image %s not found\n", r.URL.Path)
		return
	}
	content, err := s.Provider.ImageContent(r.URL.Path)
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeContent(w, r, image.Name, image.Time, content)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if IsImage(r.URL.Path) {
		s.image(w, r)
	} else {
		s.album(w, r)
	}
}
