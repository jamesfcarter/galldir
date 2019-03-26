package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jamesfcarter/galldir"
	"github.com/jamesfcarter/galldir/data"
)

func main() {
	dir := flag.String("dir", "", "Directory to serve")
	addr := flag.String("addr", "", "Address to serve")
	flag.Parse()

	provider := galldir.NewProvider(galldir.FsBackend(*dir))
	server := &galldir.Server{Provider: provider}

	assets := http.FileServer(data.Assets)

	for _, dir := range []string{"/img/", "/js/", "/css/", "/fonts/"} {
		http.Handle(dir, assets)
	}
	http.Handle("/", server)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
