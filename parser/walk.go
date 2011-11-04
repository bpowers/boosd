// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Visitor, list []*Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmtList(v Visitor, list []Stmt) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkDeclList(v Visitor, list []Decl) {
	for _, x := range list {
		Walk(v, x)
	}
}

// TODO(gri): Investigate if providing a closure to Walk leads to
//            simpler use (and may help eliminate Inspect in turn).

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {

	case *Field:
		walkIdentList(v, n.Names)
		Walk(v, n.Type)
		if n.Tag != nil {
			Walk(v, n.Tag)
		}

	case *FieldList:
		for _, f := range n.List {
			Walk(v, f)
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do


	case *CompositeLit:
		if n.Type != nil {
			Walk(v, n.Type)
		}
		walkExprList(v, n.Elts)

	case *ParenExpr:
		Walk(v, n.X)

	case *SelectorExpr:
		Walk(v, n.X)
		Walk(v, n.Sel)

	case *IndexExpr:
		Walk(v, n.X)
		Walk(v, n.Index)

	case *CallExpr:
		Walk(v, n.Fun)
		walkExprList(v, n.Args)

	case *UnaryExpr:
		Walk(v, n.X)

	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)

	case *KeyValueExpr:
		Walk(v, n.Key)
		Walk(v, n.Value)

	case *InterfaceType:
		Walk(v, n.Methods)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		Walk(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *ExprStmt:
		Walk(v, n.X)

	case *AssignStmt:
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)

	case *BlockStmt:
		walkStmtList(v, n.List)

	// Declarations
	case *ImportSpec:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		Walk(v, n.Path)

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		for _, s := range n.Specs {
			Walk(v, s)
		}

	case *InterfaceDecl:
		Walk(v, n.Name)
		Walk(v, n.Super)
		Walk(v, n.Units)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ModelDecl:
		if n.Recv != nil {
			Walk(v, n.Recv)
		}
		Walk(v, n.Name)
		Walk(v, n.Super)
		Walk(v, n.Units)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	// Files and packages
	case *File:
		Walk(v, n.Name)
		walkDeclList(v, n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for _, f := range n.Files {
			Walk(v, f)
		}

	default:
		fmt.Printf("ast.Walk: unexpected node type %T", n)
		panic("ast.Walk")
	}

	v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// for all the non-nil children of node, recursively.
//
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
