package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/jamesfcarter/galldir"
	"github.com/jamesfcarter/galldir/data"
	"github.com/jamesfcarter/galldir/s3"
)

func filesystem(dir string) http.FileSystem {
	if strings.HasPrefix(dir, "https://s3.") {
		return s3.New(dir)
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
