package main

import (
	. "boosd/parser"
	"fmt"
	"go/ast"
)

type goAST struct{}

func (p *goAST) Inspect(node Node) {
	Inspect(node, func(n Node) bool { return p.Visit(n) })
}

func (p *goAST) Visit(node Node) bool {
	// inspect is called with nil at the end of a production
	if node == nil {
		return true
	}

	switch n := node.(type) {
	case *File:
		//typeScope = append(typeScope, n.Scope)
		fmt.Println("appending typeScope")
	case *ModelDecl:
		//objScope = append(objScope, n.Objects)
		fmt.Println("model (appending objScope)", n.Name.Name)
	case *DeclStmt:
	case *AssignStmt:
	case *RefExpr:
	}

	return true
}

func passGoAST(f *File) (*ast.File, error) {
	pass := goAST{}
	pass.Inspect(f)
	return nil, fmt.Errorf("not implemented")
}
