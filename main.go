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
//go:embed core-modern.lisp

var fs embed.FS

func main() {
	file, err := fs.Open("core-modern.lisp")
	if err != nil {
		log.Fatalln(err)
	}
	exprs, err := parser.Parse("core-modern.lisp", file)
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range exprs {
		result, err := e.Eval(core.BuiltinScope())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.Print())
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
