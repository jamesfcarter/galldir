package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jamesfcarter/galldir"
)

func main() {
	dir := flag.String("dir", "", "Directory to serve")
	addr := flag.String("addr", "", "Address to serve")
	flag.Parse()

	provider := galldir.Provider{galldir.FsBackend(*dir)}
	server := &galldir.Server{Provider: provider}

	log.Fatal(http.ListenAndServe(*addr, server))
}
