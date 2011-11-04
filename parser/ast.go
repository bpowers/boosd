package parser

import (
	"boosd/token"
)

type ObjectKind int

const (
	ObjModel ObjectKind = iota
	ObjInterface
	ObjKynd
	ObjFlow
	ObjStock
	ObjTime
	ObjAux
	ObjString
)

// ----------------------------------------------------------------------------
// Interfaces
//
// There are 3 main classes of nodes: Expressions and type nodes,
// statement nodes, and declaration nodes. The node names usually
// match the corresponding Go spec production names to which they
// correspond. The node fields correspond to the individual parts
// of the respective productions.
//
// All nodes contain position information marking the beginning of
// the corresponding source text segment; it is accessible via the
// Pos accessor method. Nodes may contain additional position info
// for language constructs where comments may be found between parts
// of the construct (typically any larger, parenthesized subpart).
// That position information is needed to properly position comments
// when printing the construct.

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
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
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		Obj     *Object   // denoted object; or nil
	}
	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A CompositeLit node represents a composite literal.
	CompositeLit struct {
		Type   Expr      // literal type; or nil
		Lbrace token.Pos // position of "{"
		Elts   []Expr    // list of composite elements; or nil
		Rbrace token.Pos // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen token.Pos // position of "("
		X      Expr      // parenthesized expression
		Rparen token.Pos // position of ")"
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Lbrack token.Pos // position of "["
		Index  Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun    Expr      // function expression
		Lparen token.Pos // position of "("
		Args   []Expr    // function arguments; or nil
		Rparen token.Pos // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	//
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A PairExpr node represents a pair in a table expression.
	PairExpr struct {
		X Expr // left operand
		Y Expr // right operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	UnitExpr struct {
		X    Expr // the expression
		Unit Expr // the expression's units
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	//
	KeyValueExpr struct {
		Key   Expr
		Colon token.Pos // position of ":"
		Value Expr
	}

	// A CallExpr node represents an expression followed by an argument list.
	TableExpr struct {
		Lbrack token.Pos
		Pairs  []Expr // function arguments; or nil
		Rbrack token.Pos
	}
)

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
//
type (
	// A ModelType node represents a model type.
	ModelType struct {
		Model      token.Pos  // position of "struct" keyword
		Fields     *FieldList // list of field declarations
		Incomplete bool       // true if (source) fields are missing in the Fields list
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  token.Pos  // position of "interface" keyword
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
	Opening token.Pos // position of opening parenthesis/brace, if any
	List    []*Field  // field list; or nil
	Closing token.Pos // position of closing parenthesis/brace, if any
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
func (x *PairExpr) Pos() token.Pos      { return x.X.Pos() }
func (x *TableExpr) Pos() token.Pos     { return x.Lbrack }
func (x *UnitExpr) Pos() token.Pos      { return x.X.Pos() }
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
func (x *TableExpr) End() token.Pos     { return x.Rbrack + 1 }
func (x *PairExpr) End() token.Pos      { return x.Y.End() }
func (x *UnitExpr) End() token.Pos      { return x.Unit.End() }
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
func (x *TableExpr) exprNode()    {}
func (x *PairExpr) exprNode()     {}
func (x *UnitExpr) exprNode()     {}
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

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
//
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately preceding semicolon.
	//
	EmptyStmt struct {
		Semicolon token.Pos // position of preceding ";"
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	//
	ExprStmt struct {
		X Expr // expression
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	//
	AssignStmt struct {
		Lhs    Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // assignment token, DEFINE
		Rhs    Expr
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}"
	}
)

// Pos and End implementations for statement nodes.
//
func (s *BadStmt) Pos() token.Pos    { return s.From }
func (s *DeclStmt) Pos() token.Pos   { return s.Decl.Pos() }
func (s *EmptyStmt) Pos() token.Pos  { return s.Semicolon }
func (s *ExprStmt) Pos() token.Pos   { return s.X.Pos() }
func (s *AssignStmt) Pos() token.Pos { return s.Lhs.Pos() }
func (s *BlockStmt) Pos() token.Pos  { return s.Lbrace }

func (s *BadStmt) End() token.Pos  { return s.To }
func (s *DeclStmt) End() token.Pos { return s.Decl.End() }
func (s *EmptyStmt) End() token.Pos {
	return s.Semicolon + 1 /* len(";") */
}
func (s *ExprStmt) End() token.Pos   { return s.X.End() }
func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }
func (s *BlockStmt) End() token.Pos  { return s.Rbrace + 1 }

// stmtNode() ensures that only statement nodes can be
// assigned to a StmtNode.
//
func (s *BadStmt) stmtNode()    {}
func (s *DeclStmt) stmtNode()   {}
func (s *EmptyStmt) stmtNode()  {}
func (s *ExprStmt) stmtNode()   {}
func (s *AssignStmt) stmtNode() {}
func (s *BlockStmt) stmtNode()  {}

// ----------------------------------------------------------------------------
// Declarations

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
//
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Name   *Ident    // local package name (including "."); or nil
		Path   *BasicLit // import path
		EndPos token.Pos // end of spec (overrides Path.Pos if nonzero)
	}

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	//
	KindSpec struct {
		Names  []*Ident // value names (len(Names) > 0)
		Type   Expr     // value type; or nil
		Values []Expr   // initial values; or nil
	}
)

// Pos and End implementations for spec nodes.
//
func (s *ImportSpec) Pos() token.Pos {
	if s.Name != nil {
		return s.Name.Pos()
	}
	return s.Path.Pos()
}
func (s *KindSpec) Pos() token.Pos { return s.Names[0].Pos() }

func (s *ImportSpec) End() token.Pos {
	if s.EndPos != 0 {
		return s.EndPos
	}
	return s.Path.End()
}

func (s *KindSpec) End() token.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}

// specNode() ensures that only spec nodes can be
// assigned to a Spec.
//
func (s *ImportSpec) specNode() {}
func (s *KindSpec) specNode()   {}

// A declaration is represented by one of the following declaration nodes.
//
type (
	// A BadDecl node is a placeholder for declarations containing
	// syntax errors for which no correct declaration nodes can be
	// created.
	//
	BadDecl struct {
		From, To token.Pos // position range of bad declaration
	}

	// A GenDecl node (generic declaration node) represents an import,
	// constant, type or variable declaration. A valid Lparen position
	// (Lparen.Line > 0) indicates a parenthesized declaration.
	//
	// Relationship between Tok value and Specs element type:
	//
	//	token.IMPORT  *ImportSpec
	//	token.CONST   *ValueSpec
	//	token.TYPE    *TypeSpec
	//	token.VAR     *ValueSpec
	//
	GenDecl struct {
		TokPos token.Pos   // position of Tok
		Tok    token.Token // IMPORT, CONST, TYPE, VAR
		Lparen token.Pos   // position of '(', if any
		Specs  []Spec
		Rparen token.Pos // position of ')', if any
	}

	VarDecl struct {
		Name  *Ident // name of the variable
		Type  *Ident // type (stock, flow) of the variable
		Units Expr   // name of the variable
	}

	// A FuncDecl node represents a function declaration.
	InterfaceDecl struct {
		Recv *FieldList // receiver (methods); or nil (functions)
		Name *Ident     // function/method name
		Type *ModelType // position of Func keyword, parameters and results
		Body *BlockStmt // function body; or nil (forward declaration)
	}

	// A FuncDecl node represents a function declaration.
	ModelDecl struct {
		Recv *FieldList // receiver (methods); or nil (functions)
		Name *Ident     // function/method name
		Type *ModelType // position of Func keyword, parameters and results
		Body *BlockStmt // function body; or nil (forward declaration)
	}
)

// Pos and End implementations for declaration nodes.
//
func (d *BadDecl) Pos() token.Pos       { return d.From }
func (d *GenDecl) Pos() token.Pos       { return d.TokPos }
func (d *VarDecl) Pos() token.Pos       { return d.Name.Pos() }
func (d *InterfaceDecl) Pos() token.Pos { return d.Type.Pos() }
func (d *ModelDecl) Pos() token.Pos     { return d.Type.Pos() }

func (d *BadDecl) End() token.Pos { return d.To }
func (d *VarDecl) End() token.Pos { return d.Name.End() } // FIXME
func (d *GenDecl) End() token.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Specs[0].End()
}
func (d *InterfaceDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}
func (d *ModelDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}

// declNode() ensures that only declaration nodes can be
// assigned to a DeclNode.
//
func (d *BadDecl) declNode()       {}
func (d *GenDecl) declNode()       {}
func (d *VarDecl) declNode()       {}
func (d *InterfaceDecl) declNode() {}
func (d *ModelDecl) declNode()     {}

// ----------------------------------------------------------------------------
// Files and packages

// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.
//
type File struct {
	Package    token.Pos     // position of "package" keyword
	Name       *Ident        // package name
	Decls      []Decl        // top-level declarations; or nil
	Scope      *Scope        // package scope (this file only)
	Imports    []*ImportSpec // imports in this file
	Unresolved []*Ident      // unresolved identifiers in this file
}

func (f *File) Pos() token.Pos { return f.Package }
func (f *File) End() token.Pos {
	if n := len(f.Decls); n > 0 {
		return f.Decls[n-1].End()
	}
	return f.Name.End()
}

// A Package node represents a set of source files
// collectively building a Go package.
//
type Package struct {
	Name    string             // package name
	Scope   *Scope             // package scope across all files
	Imports map[string]*Object // map of package id -> package object
	Files   map[string]*File   // Go source files by filename
}

func (p *Package) Pos() token.Pos { return token.NoPos }
func (p *Package) End() token.Pos { return token.NoPos }
