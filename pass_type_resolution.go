package main

import (
	. "boosd/parser"
	"boosd/token"
	"fmt"
	"strings"
)

func x() {

}

type Pass interface {
	Inspect(Node) bool
}

type typeResolution struct {
	fset      *token.FileSet
	typeScope *Scope
	objScope  *Scope
	currModel *ModelDecl
	depth     int
}

func (p *typeResolution) Inspect(node Node) {
	Inspect(node, func(n Node) bool { return p.Visit(n) })
}

func (p *typeResolution) Visit(node Node) bool {
	// inspect is called with nil at the end of a production
	if node == nil {
		p.depth--
		return true
	}

	switch n := node.(type) {
	case *File:
		//typeScope = append(typeScope, n.Scope)
		fmt.Println("appending typeScope")
	case *ModelDecl:
		p.currModel = n
		//objScope = append(objScope, n.Objects)
		fmt.Println("model (appending objScope)", n.Name.Name)
	case *DeclStmt:
		p.currModel.Virtual = true
		fmt.Println(strings.Repeat("  ", p.depth) + n.Name())
	case *AssignStmt:
		fmt.Println(strings.Repeat("  ", p.depth) + n.Name())
	case *RefExpr:
		if p.objScope != nil && p.objScope.Lookup(n.Name) == nil {
			panic("unknown variable referenced: " + n.Name)
		}
	}
	fmt.Printf("%s(%T)\n", strings.Repeat("  ", p.depth), node)
	//fmt.Printf("%s%#v\n", indent, node)
	p.depth++

	return true
}

func passTypeResolution(f *File) {

	pass := typeResolution{}

	pass.Inspect(f)
}
