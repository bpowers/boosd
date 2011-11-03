package main

import (
	"boosd/token"
)

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

type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
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

type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos     // identifier position
		Name    string  // identifier name
		Obj     *Object // denoted object; or nil
	}
	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos    // literal position
		Kind     token.Token  // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A CompositeLit node represents a composite literal.
	CompositeLit struct {
		Type   Expr   // literal type; or nil
		Lbrace token.Pos    // position of "{"
		Elts   []Expr // list of composite elements; or nil
		Rbrace token.Pos    // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen token.Pos  // position of "("
		X      Expr // parenthesized expression
		Rparen token.Pos  // position of ")"
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr // expression
		Lbrack token.Pos  // position of "["
		Index  Expr // index expression
		Rbrack token.Pos  // position of "]"
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun      Expr   // function expression
		Lparen   token.Pos    // position of "("
		Args     []Expr // function arguments; or nil
		Ellipsis token.Pos    // position of "...", if any
		Rparen   token.Pos    // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	//
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr  // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr  // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr  // right operand
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	//
	KeyValueExpr struct {
		Key   Expr
		Colon token.Pos // position of ":"
		Value Expr
	}
)

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
//
type (
	// A ModelType node represents a model type.
	ModelType struct {
		Model      token.Pos        // position of "struct" keyword
		Fields     *FieldList // list of field declarations
		Incomplete bool       // true if (source) fields are missing in the Fields list
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  token.Pos        // position of "interface" keyword
		Methods    *FieldList // list of methods
		Incomplete bool       // true if (source) methods are missing in the Methods list
	}
)

// ----------------------------------------------------------------------------
// Expressions and types

// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.
//
type Field struct {
	Names []*Ident  // field/method/parameter names; or nil if anonymous field
	Type  Expr      // field/method/parameter type
	Value Expr      // field/method/parameter value
	Tag   *BasicLit // field tag; or nil
}

func (f *Field) Pos() token.Pos {
	if len(f.Names) > 0 {
		return f.Names[0].Pos()
	}
	return f.Type.Pos()
}

func (f *Field) End() token.Pos {
	if f.Tag != nil {
		return f.Tag.End()
	}
	return f.Type.End()
}

// A FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList struct {
	Opening token.Pos      // position of opening parenthesis/brace, if any
	List    []*Field // field list; or nil
	Closing token.Pos      // position of closing parenthesis/brace, if any
}

func (f *FieldList) Pos() token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return token.NoPos
}

func (f *FieldList) End() token.Pos {
	if f.Closing.IsValid() {
		return f.Closing + 1
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return token.NoPos
}

func (x *Ident) Pos() token.Pos    { return x.NamePos }
func (x *BasicLit) Pos() token.Pos { return x.ValuePos }
func (x *CompositeLit) Pos() token.Pos {
	if x.Type != nil {
		return x.Type.Pos()
	}
	return x.Lbrace
}
func (x *ParenExpr) Pos() token.Pos     { return x.Lparen }
func (x *SelectorExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *IndexExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos      { return x.Fun.Pos() }
func (x *UnaryExpr) Pos() token.Pos     { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos    { return x.X.Pos() }
func (x *KeyValueExpr) Pos() token.Pos  { return x.Key.Pos() }
func (x *ModelType) Pos() token.Pos     { return x.Model }
func (x *InterfaceType) Pos() token.Pos { return x.Interface }

func (x *Ident) End() token.Pos         { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos      { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *CompositeLit) End() token.Pos  { return x.Rbrace + 1 }
func (x *ParenExpr) End() token.Pos     { return x.Rparen + 1 }
func (x *SelectorExpr) End() token.Pos  { return x.Sel.End() }
func (x *IndexExpr) End() token.Pos     { return x.Rbrack + 1 }
func (x *CallExpr) End() token.Pos      { return x.Rparen + 1 }
func (x *UnaryExpr) End() token.Pos     { return x.X.End() }
func (x *BinaryExpr) End() token.Pos    { return x.Y.End() }
func (x *KeyValueExpr) End() token.Pos  { return x.Value.End() }
func (x *ModelType) End() token.Pos     { return x.Fields.End() }
func (x *InterfaceType) End() token.Pos { return x.Methods.End() }

// exprNode() ensures that only expression/type nodes can be
// assigned to an ExprNode.
//
func (x *Ident) exprNode()        {}
func (x *BasicLit) exprNode()     {}
func (x *CompositeLit) exprNode() {}
func (x *ParenExpr) exprNode()    {}
func (x *SelectorExpr) exprNode() {}
func (x *IndexExpr) exprNode()    {}
func (x *CallExpr) exprNode()     {}
func (x *UnaryExpr) exprNode()    {}
func (x *BinaryExpr) exprNode()   {}
func (x *KeyValueExpr) exprNode() {}

func (x *ModelType) exprNode()     {}
func (x *InterfaceType) exprNode() {}

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

// ----------------------------------------------------------------------------
// Convenience functions for Idents

var noPos token.Pos

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
//
func NewIdent(name string) *Ident { return &Ident{noPos, name, nil} }
