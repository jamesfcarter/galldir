// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	fs := http.Dir("assets")
	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:     "data",
		VariableName:    "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
