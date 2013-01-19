package main

import (
	. "boosd/parser"
	"fmt"
	"strings"
)

type timeSpec struct {
	currModel *ModelDecl
	scope     *Scope
	depth     int
}

func (p *timeSpec) Inspect(node Node) {
	Inspect(node, func(n Node) bool { return p.Visit(n) })
}

func (p *timeSpec) Visit(node Node) bool {
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
		if p.scope != nil && p.scope.Lookup(n.Name) == nil {
			panic("unknown variable referenced: " + n.Name)
		}
	}
	fmt.Printf("%s(%T)\n", strings.Repeat("  ", p.depth), node)
	//fmt.Printf("%s%#v\n", indent, node)
	p.depth++

	return true
}

func passTimespec(f *File) {
	pass := timeSpec{scope: f.Scope}
	pass.Inspect(f)
}
