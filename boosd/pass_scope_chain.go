// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boosd

import (
	"fmt"
	"os"
	"strings"
)

type scopeChain struct {
	currModel *ModelDecl
	scope     *Scope
	depth     int
}

func (p *scopeChain) Inspect(node Node) {
	Inspect(node, func(n Node) bool { return p.Visit(n) })
}

func (p *scopeChain) Visit(node Node) bool {
	// inspect is called with nil at the end of a production
	if node == nil {
		p.depth--
		return true
	}

	//fmt.Fprintf(os.Stderr, "%s(%T)\n", strings.Repeat("  ", p.depth), node)
	switch n := node.(type) {
	case *File:
		//typeScope = append(typeScope, n.Scope)
		//fmt.Println("appending typeScope")
	case *ModelDecl:
		p.currModel = n
		//objScope = append(objScope, n.Objects)
		//fmt.Println("model (appending objScope)", n.Name.Name)
	case *DeclStmt:
		p.currModel.Virtual = true
		//fmt.Println(strings.Repeat("  ", p.depth) + n.Name())
	case *AssignStmt:
		//fmt.Println(strings.Repeat("  ", p.depth) + n.Name())
	case *RefExpr:
		if p.scope != nil && p.scope.Lookup(n.Name) == nil {
			panic("unknown variable referenced: " + n.Name)
		}
	}
	//fmt.Printf("%s%#v\n", indent, node)
	p.depth++

	return true
}

func PassScopeChain(f *File) {
	pass := scopeChain{scope: f.Scope}
	pass.Inspect(f)
}
