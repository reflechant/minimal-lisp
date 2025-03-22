package parser

import (
	"slices"
	"strings"
	"testing"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyList(t *testing.T) {
	rdr := strings.NewReader("()")
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewList(1, 1),
	}
	assert.Equal(t, expected, exprs)
}

func TestOneAtom(t *testing.T) {
	rdr := strings.NewReader("foo")
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewSymbol(1, 1, "foo"),
	}
	assert.Equal(t, expected, exprs)
}

func TestNestedList(t *testing.T) {
	rdr := strings.NewReader("(foo ( bar) baz)")
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewList(1, 1,
			core.NewSymbol(1, 2, "foo"),
			core.NewList(1, 6, core.NewSymbol(1, 8, "bar")),
			core.NewSymbol(1, 13, "baz"),
		),
	}
	for i := range expected {
		assert.Equal(t, expected[i].String(), exprs[i].String())
	}
}

func TestMultiLine(t *testing.T) {
	rdr := strings.NewReader(
		`()
  foo
()`)
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewList(1, 1),
		core.NewSymbol(2, 3, "foo"),
		core.NewList(3, 1),
	}
	for i := range expected {
		assert.Equal(t, expected[i].String(), exprs[i].String())
	}
}

func TestSyntaxErrors(t *testing.T) {
	cases := []struct {
		input          string
		expectedErrMsg string
	}{
		{
			input:          "(",
			expectedErrMsg: "list opened at 1:1 was not closed",
		},
		{
			input:          "(quote 2",
			expectedErrMsg: "list opened at 1:1 was not closed",
		},
		{
			input:          ")",
			expectedErrMsg: "unexpected token )",
		},
		{
			input:          "quote 2)",
			expectedErrMsg: "unexpected token )",
		},
		{
			input:          "'",
			expectedErrMsg: "unexpected end of input: quote needs an argument",
		},
	}

	for tc := range slices.Values(cases) {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			rdr := strings.NewReader(tc.input)
			_, err := Parse("test", rdr)
			require.ErrorContains(t, err, tc.expectedErrMsg)

		})
	}
}
