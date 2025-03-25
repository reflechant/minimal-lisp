package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/reflechant/minimal-lisp/parser"
)

const prompt = ">>> "

func REPL(scope core.Scope, in io.Reader, out io.Writer) error {
	// print the REPL prompt
	_, err := out.Write([]byte(prompt))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		rdr := strings.NewReader(line)
		exprs, err := parser.Parse("repl", rdr)
		if err != nil {
			_, err := out.Write(fmt.Appendf(nil, "error parsing REPL input: %v\n", err))
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		}

		for _, e := range exprs {
			result, err := e.Eval(scope)
			if err != nil {
				out.Write(fmt.Appendf(nil, "error evaluating REPL input: %v\n", err))
				continue
			}

			if result != nil {
				out.Write([]byte(result.String()))
				out.Write([]byte{'\n'})
				continue
			}
			out.Write([]byte("nil\n"))
		}

		// print the REPL prompt
		_, err = out.Write([]byte(prompt))
		if err != nil {
			return err
		}
	}

	return nil
}
