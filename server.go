package galldir

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"time"
)

// Server implements to http.Handler interface to serve a photo gallery
type Server struct {
	Provider *Provider
	Assets   http.FileSystem
}

const (
	albumPath  = "/img/album.png"
	thumbWidth = 250
)

func (s *Server) albumThumb(w http.ResponseWriter, r *http.Request, album *Album) {
	content, err := s.Provider.CoverThumb(album, thumbWidth)
	if err != nil {
		log.Println(err)
		content, err = s.assetThumb(albumPath)
	}
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeContent(w, r, "", time.Now(), content)
}

func (s *Server) album(w http.ResponseWriter, r *http.Request) {
	refresh := cacheRefresh(r)
	album, err := s.Provider.Album(r.URL.Path, refresh)
	if err != nil {
		log.Println(err)
		return
	}
	if isThumb(r) {
		s.albumThumb(w, r, album)
		return
	}
	page := struct {
		Refresh template.URL
		Album   *Album
	}{
		Refresh: func() template.URL {
			if refresh {
				return template.URL("&refresh=1")
			}
			return template.URL("")
		}(),
		Album: album,
	}
	err = indexTemplate.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

func requestFlag(r *http.Request, flag string) bool {
	thumbParams, ok := r.URL.Query()[flag]
	return ok && len(thumbParams[0]) > 0
}

func cacheRefresh(r *http.Request) bool {
	return requestFlag(r, "refresh")
}

func isThumb(r *http.Request) bool {
	return requestFlag(r, "thumb")
}

func (s *Server) assetThumb(path string) (io.ReadSeeker, error) {
	cacheName := ThumbName("assetthumb", thumbWidth, path)
	image, err := s.Assets.Open(path)
	if err != nil {
		return nil, err
	}
	return s.Provider.CachedThumb(cacheName, thumbWidth, image)
}

func (s *Server) image(w http.ResponseWriter, r *http.Request) {
	albumPath := path.Dir(r.URL.Path)
	album, err := s.Provider.Album(albumPath, false)
	if err != nil {
		log.Printf("Failed to fetch album %s: %v\n", albumPath, err)
		return
	}
	image := album.Image(r.URL.Path)
	if image == nil {
		log.Printf("image %s not found\n", r.URL.Path)
		return
	}
	var content io.ReadSeeker
	if isThumb(r) {
		content, err = s.Provider.ImageThumb(r.URL.Path, thumbWidth)
		if err != nil {
			content, err = s.assetThumb(albumPath)
		}
	} else {
		content, err = s.Provider.ImageContent(r.URL.Path)
	}
	if err != nil {
		log.Printf("failed to serve image %s: %v\n", r.URL.Path, err)
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
	<title>{{ .Album.Name }}</title>
	<link type="text/css" rel="stylesheet" href="/css/lightgallery.css" />
	<link type="text/css" rel="stylesheet" href="/css/galldir.css" />
    </head>
    <body>
	<h1>{{ .Album.Name }}</h1>
        <script src="/js/lightgallery.min.js"></script>
        <script src="/js/lg-thumbnail.min.js"></script>
        <script src="/js/lg-fullscreen.min.js"></script>
	<div class="galldir-albums">
	    {{ range .Album.Albums }}
		<figure><p><a href="{{ .Path }}">
			<img src="{{ .Path }}?thumb=1{{ $.Refresh }}" />
			<figcaption>{{ .Name }}</figcaption>
		</a></p></figure>
	    {{ end }}
	</div>
	<div id="lightgallery">
	{{ range .Album.Photos }}
	    <a href="{{ .Path }}"><img src="{{ .Path }}?thumb=1" /></a>
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
