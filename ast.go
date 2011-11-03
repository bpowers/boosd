package main

import ()

type ObjKind int

const (
	ObjModel ObjKind = iota
	ObjInterface
	ObjKynd
	ObjFlow
	ObjStock
	ObjTime
	ObjAux
	ObjString
)

type Pos struct {
	Line int
	Col  int
}

type Node interface {
	Pos() Pos // position of first character belonging to the node
	End() Pos // position of first character immediately after the node
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

type Object struct {
	Name string
	Kind ObjKind
	Decl interface{}
	Data interface{}
	Type string
	Unit string
}
type Token int

type (
	BinaryExpr struct {
		X     Expr  // left operand
		OpPos Pos   // position of Op
		Op    Token // operator
		Y     Expr  // right operand
	}
)

func (x *BinaryExpr) Pos() Pos { return x.X.Pos() }
func (x *BinaryExpr) End() Pos { return x.Y.End() }

// exprNode() ensures that only expression/type nodes can be
// assigned to an ExprNode.
//
func (x *BinaryExpr) exprNode() {}

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
// based off Go's AST walk in pkg/go/ast
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	switch n := node.(type) {
	case *BinaryExpr:
		Walk(v, n.X)
	}

	v.Visit(nil)
}
