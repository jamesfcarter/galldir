package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jamesfcarter/galldir"
	"github.com/jamesfcarter/galldir/data"
	s3 "github.com/jamesfcarter/s3httpfilesystem"
)

func s3ConfigFromURL(uri string) (endpoint, region, bucket string) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	endpoint = u.Host
	bucket = strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)[0]
	hostBits := strings.SplitN(endpoint, ".", 3)
	if len(hostBits) > 1 {
		region = hostBits[1]
	}
	return
}

func filesystem(dir string) http.FileSystem {
	if strings.HasPrefix(dir, "https://s3.") {
		endpoint, region, bucket := s3ConfigFromURL(dir)
		return s3.New(endpoint, region, bucket)
	}
	return http.Dir(dir)
}

func main() {
	dir := flag.String("dir", "", "Directory to serve")
	addr := flag.String("addr", "", "Address to serve")
	flag.Parse()

	provider := galldir.NewProvider(filesystem(*dir))
	server := &galldir.Server{
		Provider: provider,
		Assets:   data.Assets,
	}

	assets := http.FileServer(data.Assets)

	for _, dir := range []string{
		"/favicon.ico", "/img/", "/js/", "/css/", "/fonts/",
	} {
		http.Handle(dir, assets)
	}
	http.Handle("/", server)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
