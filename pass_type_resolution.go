package main

import (
	. "boosd/parser"
	"fmt"
	"strings"
)

func passTypeResolution(f *File) {
	var curr Node
	typeScope := make([]*Scope, 2)
	objScope := make([]*Scope, 2)
	depth := 0
	var currModel *ModelDecl

	Inspect(f, func(node Node) bool {
		// inspect is called with nil at the end of a production
		if node == nil {
			depth--
			return true
		}

		curr = node

		switch n := node.(type) {
		case *File:
			typeScope = append(typeScope, n.Scope)
			fmt.Println("appending typeScope")
		case *ModelDecl:
			currModel = n
			objScope = append(objScope, n.Objects)
			fmt.Println("model (appending objScope)", n.Name.Name)
		case *DeclStmt:
			currModel.Virtual = true
			fmt.Println(strings.Repeat("  ", depth) + n.Name())
		case *AssignStmt:
			fmt.Println(strings.Repeat("  ", depth) + n.Name())
		}
		fmt.Printf("%s%d (%T)\n", strings.Repeat("  ", depth), depth, node)
		//fmt.Printf("%s%#v\n", indent, node)
		depth++

		return true
	})
}
