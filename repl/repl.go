package repl

import (
	"bufio"
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
			_, err := out.Write(fmt.Appendf(nil, "%v\n", err))
			if err != nil {
				return err
			}
		}

		for _, e := range exprs {
			result, err := e.Eval(scope)
			if err != nil {
				_, err := out.Write(fmt.Appendf(nil, "%v\n", err))
				if err != nil {
					return err
				}
				continue
			}

			if result != nil {
				_, err := out.Write(fmt.Appendf(nil, "%v\n", result.String()))
				if err != nil {
					return err
				}
				continue
			}
		}

		// print the REPL prompt
		_, err = out.Write([]byte(prompt))
		if err != nil {
			return err
		}
	}

	return nil
}
