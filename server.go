package galldir

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Server implements to http.Handler interface to serve a photo gallery
type Server struct {
	Provider *Provider
	Assets   http.FileSystem
}

const (
	albumPath = "/img/album.png"
)

func (s *Server) albumThumb(w http.ResponseWriter, r *http.Request, album *Album, thumbSize int) {
	content, err := s.Provider.CoverThumb(album, thumbSize)
	if err != nil {
		log.Println(err)
		content, err = s.assetThumb(albumPath, thumbSize)
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
	thumbSize, needThumb := isThumb(r)
	if needThumb {
		s.albumThumb(w, r, album, thumbSize)
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

func requestParamInt(r *http.Request, flag string) (int, bool) {
	thumbParams := r.URL.Query()
	valueString := thumbParams.Get(flag)
	if valueString == "" {
		return 0, false
	}
	value, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, false
	}
	return value, true
}

func cacheRefresh(r *http.Request) bool {
	_, ok := requestParamInt(r, "refresh")
	return ok
}

func isThumb(r *http.Request) (int, bool) {
	return requestParamInt(r, "thumb")
}

func (s *Server) assetThumb(path string, thumbSize int) (io.ReadSeeker, error) {
	cacheName := ThumbName("assetthumb", thumbSize, path)
	image, err := s.Assets.Open(path)
	if err != nil {
		return nil, err
	}
	return s.Provider.CachedThumb(cacheName, thumbSize, image)
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
	thumbSize, needThumb := isThumb(r)
	if needThumb {
		content, err = s.Provider.ImageThumb(r.URL.Path, thumbSize)
		if err != nil {
			content, err = s.assetThumb(albumPath, thumbSize)
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
			<img src="{{ .Path }}?thumb=250{{ $.Refresh }}" />
			<figcaption>{{ .Name }}</figcaption>
		</a></p></figure>
	    {{ end }}
	</div>
	<div id="lightgallery">
	{{ range .Album.Photos }}
	    <a href="{{ .Path }}"><img src="{{ .Path }}?thumb=250" /></a>
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
