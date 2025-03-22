package main_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/reflechant/minimal-lisp/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	cases := []struct {
		input          string
		expected       string
		expectedErrMsg string
	}{
		{
			input:    "()",
			expected: "()",
		},
		{
			input:          "x",
			expectedErrMsg: "unbound symbol",
		},
		{
			input:          "(f 2)",
			expectedErrMsg: "evaluation error",
		},
		// built-in functions
		// quote
		{
			input:    "'s",
			expected: "s",
		},
		{
			input:    "'(1 2)",
			expected: "(1 2)",
		},
		{
			input:    "'(1 ())",
			expected: "(1 ())",
		},
		// atom
		{
			input:    "(atom 'x)",
			expected: "t",
		},
		{
			input:    "(atom '())",
			expected: "t",
		},
		{
			input:    "(atom '(a b c))",
			expected: "()",
		},
		{
			input:          "(atom)",
			expectedErrMsg: "atom: expects 1 argument, 0 given",
		},
		// eq
		{
			input:    "(eq 'a 'a)",
			expected: "t",
		},
		{
			input:    "(eq 'a 'b)",
			expected: "()",
		},
		{
			input:    "(eq '() '())",
			expected: "t",
		},
		{
			input:    "(eq '(1) '(1))",
			expected: "()",
		},
		{
			input:          "(eq 'a)",
			expectedErrMsg: "eq: expects 2 arguments, got 1",
		},
		{
			input:          "(eq)",
			expectedErrMsg: "eq: expects 2 arguments, got 0",
		},
		// car
		{
			input:    "(car '(a b c))",
			expected: "a",
		},
		{
			input:    "(car '())",
			expected: "()",
		},
		{
			input:          "(car 'x)",
			expectedErrMsg: "car: argument must be a list",
		},
		{
			input:          "(car)",
			expectedErrMsg: "car: expects 1 argument, got 0",
		},
		// cdr
		{
			input:    "(cdr '(a b c))",
			expected: "(b c)",
		},
		{
			input:    "(cdr '())",
			expected: "()",
		},
		{
			input:          "(cdr 'x)",
			expectedErrMsg: "cdr: argument must be a list",
		},
		{
			input:          "(cdr)",
			expectedErrMsg: "cdr: expects 1 argument, got 0",
		},
		// cons
		{
			input:    "(cons 'a '(b c))",
			expected: "(a b c)",
		},
		{
			input:    "(cons 'a (cons 'b (cons 'c ())))",
			expected: "(a b c)",
		},
		{
			input:          "(cons 'a 'b)",
			expectedErrMsg: "cons: 2nd argument must be a list",
		},
		{
			input:          "(cons '1)",
			expectedErrMsg: "cons: expects 2 arguments, got 1",
		},
		{
			input:          "(cons)",
			expectedErrMsg: "cons: expects 2 arguments, got 0",
		},
		// cond
		{
			input:    "(cond ((eq 'a 'a) 'b))",
			expected: "b",
		},
		// {
		// 	input:    "(cond 'x)",
		// 	expected: "()",
		// },
		// {
		// 	input:          "(cond cond)",
		// 	expectedErrMsg: "cond: argument #1 is not a list",
		// },

		// {
		// 	input:    "(cond)",
		// 	expected: "()",
		// },
	}

	for tc := range slices.Values(cases) {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			rdr := strings.NewReader(tc.input)
			exprs, err := parser.Parse("test", rdr)
			require.NoError(t, err)
			assert.Len(t, exprs, 1)

			scope := core.BuiltinScope()
			result, err := exprs[0].Eval(scope)
			if tc.expectedErrMsg != "" {
				require.ErrorContains(t, err, tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result.String())
			}
		})
	}
}
