package galldir

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
)

type Server struct {
	Provider *Provider
}

func (s *Server) album(w http.ResponseWriter, r *http.Request) {
	album, err := s.Provider.Album(r.URL.Path)
	if err != nil {
		log.Println(err)
		return
	}
	err = indexTemplate.Execute(w, album)
	if err != nil {
		log.Println(err)
	}
}

func isThumb(r *http.Request) bool {
	thumbParams, ok := r.URL.Query()["thumb"]
	return ok && len(thumbParams[0]) > 0
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
	var content io.ReadSeeker
	if isThumb(r) {
		content, err = s.Provider.ImageThumb(r.URL.Path, 100)
	} else {
		content, err = s.Provider.ImageContent(r.URL.Path)
	}
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

var indexTemplate = template.Must(template.New("index.html").Parse(`
<html>
    <head>
	<link type="text/css" rel="stylesheet" href="/css/lightgallery.css" />
    </head>
    <body>
        <script src="/js/lightgallery.min.js"></script>
        <script src="/js/lg-thumbnail.min.js"></script>
        <script src="/js/lg-fullscreen.min.js"></script>
	<div>
	    <ul>
	    {{ range .Albums }}
		<li><a href="{{ .Path }}">{{ .Name }}</a></li>
	    {{ end }}
	    </ul>
	</div>
	<div id="lightgallery">
	{{ range .Photos }}
	    <a href="{{ .Path }}">
		<img src="{{ .Path }}?thumb=1" />
	    </a>
	{{ end }}
	</div>
    	<script>
	    lightGallery(document.getElementById('lightgallery'), {
		thumbnail:true,
		animatedthumb:true
	    });
        </script>
    </body>
</html>
`))
