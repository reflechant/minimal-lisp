package parser

import (
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
		core.NewEmptyList(1, 1),
	}
	assert.Equal(t, expected, exprs)
}

func TestOneAtom(t *testing.T) {
	rdr := strings.NewReader("foo")
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewAtom(1, 1, "foo"),
	}
	assert.Equal(t, expected, exprs)
}

func TestNestedList(t *testing.T) {
	rdr := strings.NewReader("(foo ( bar) baz)")
	exprs, err := Parse("test", rdr)
	require.NoError(t, err)
	expected := []core.SExpr{
		core.NewListFromElements(1, 1, []core.SExpr{
			core.NewAtom(1, 2, "foo"),
			core.NewList(1, 6, core.NewAtom(1, 8, "bar")),
			core.NewAtom(1, 13, "baz"),
		}),
	}
	for i := range expected {
		assert.True(t, expected[i].Eq(exprs[i]), "expr %d: %v is not equal %v", i, expected[i].Print(), exprs[i].Print())
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
		core.NewEmptyList(1, 1),
		core.NewAtom(2, 3, "foo"),
		core.NewEmptyList(3, 1),
	}
	for i := range expected {
		assert.True(t, expected[i].Eq(exprs[i]), "expr %d: %v is not equal %v", i, expected[i], exprs[i])
	}
}
