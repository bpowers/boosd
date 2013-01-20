package main

import (
	"boosd/parser"
	"go/ast"
)

func genGoFile(f *parser.File) (*ast.File, error) {
	goF := &ast.File{Name: ast.NewIdent("main")}

	return goF, nil
}
