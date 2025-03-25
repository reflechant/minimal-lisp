package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/reflechant/minimal-lisp/parser"
	"github.com/reflechant/minimal-lisp/repl"
)

//go:embed core.lisp

var fs embed.FS

func Import(scope core.Scope, srcName string, src io.Reader) error {
	exprs, err := parser.Parse(srcName, src)
	if err != nil {
		return fmt.Errorf("error ", err)
	}
	for _, e := range exprs {
		_, err := e.Eval(scope)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	file, err := fs.Open("core.lisp")
	if err != nil {
		log.Fatalln(err)
	}
	scope := core.BuiltinScope()
	err = Import(scope, "core.lisp", file)
	if err != nil {
		log.Fatalln(err)
	}
	err = repl.REPL(scope, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalln(err)
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
