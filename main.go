package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/reflechant/minimal-lisp/parser"
)

//go:embed core.lisp

var fs embed.FS

func main() {
	srcName := "core.lisp"
	file, err := fs.Open(srcName)
	if err != nil {
		log.Fatalln(err)
	}
	exprs, err := parser.Parse(srcName, file)
	if err != nil {
		log.Fatalln(err)
	}

	scope := core.BuiltinScope()
	for _, e := range exprs {
		_, err := e.Eval(scope)
		if err != nil {
			log.Fatal(fmt.Errorf("%s: %w", srcName, err))
		}
	}
}

func ReadFile(fpath string) (*os.File, error) {
	fpath = path.Clean(fpath)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
