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
		PackageName:  "data",
		VariableName: "Assets",
		VariableComment: "go:generate go run assets_generate.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
