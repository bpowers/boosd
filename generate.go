package main

import (
	. "boosd/parser"
	"fmt"
	"go/ast"
)

type generator struct{}

func (g *generator) declList(list []Decl) ([]ast.Decl, error) {
	return []ast.Decl{}, nil
}

func (g *generator) gen(node Node) (interface{}, error) {
	var err error
	switch n := node.(type) {
	case *File:
		f := &ast.File{Name: ast.NewIdent("main")}
		if f.Decls, err = g.declList(n.Decls); err != nil {
			return nil, err
		}
		return f, nil
	}

	return nil, fmt.Errorf("unimplemented")
}

func generateGoAST(f *File) (*ast.File, error) {
	node, err := (&generator{}).gen(f)
	if err != nil {
		return nil, err
	}

	if goFile, ok := node.(*ast.File); ok {
		return goFile, nil
	}

	return nil, fmt.Errorf("gen(%v): node not an *ast.File", node)
}
