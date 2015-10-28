package sieve

import (
	"github.com/qingshan/sieve/parse"
	"github.com/qingshan/sieve/ast"
)

func Parse(name, text string) *ast.File {
	p := parse.Parse(name, text)
	return &ast.File{name, p.List}
}
