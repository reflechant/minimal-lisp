package main

import (
	"fmt"
	"os"
	"path"

	"github.com/reflechant/minimal-lisp/parser"
)

func main() {
	// parser.Parse()
}

func ReadFile(fpath string) (*os.File, error) {
	fpath = path.Clean(fpath)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
