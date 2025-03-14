package main

import (
	"log"
	"os"
	"path"
)

func main() {
	// err := minlisp.Repl(nil, os.Stdin, os.Stdout)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}

func ReadFile(fpath string) (*os.File, error) {
	fpath = path.Clean(fpath)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
