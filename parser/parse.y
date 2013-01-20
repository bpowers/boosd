// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

%{

package parser

import (
	"go/token"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	tok    tok
	strs   []string
	ids    []*Ident
	file   File
	id     *Ident
	str    string
	lit    *BasicLit
	exprs  []Expr
	expr   Expr
	stmt   Stmt
	tlDecl Decl
	decl   *VarDecl
	decls  []Decl
	block  *BlockStmt
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <strs>   imports
%type <ids>    id_list
%type <file>   file
%type <str>    import
%type <id>     ident specializes
%type <tok>    top_type
%type <block>  stmts
%type <stmt>   stmt
%type <expr>   expr number pair table expr_w_unit opt_kind initializer assignment ref
%type <exprs>  pairs expr_list initializers
%type <decl>   var_decl
%type <tlDecl> def
%type <decls>  defs
%type <lit>    lit

// same for terminals
%token <tok> YIMPORT YKIND YKIND_DECL YPACKAGE
%token <tok> YSPECIALIZES YINTERFACE YMODEL
%token <tok> YIDENT YLITERAL YNUMBER

%left '+'  '-'
%left '*'  '/'
%left '^'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */
%left FN_CALL

%%

file:	imports
	kinds
	defs
	{
		$$.Decls = $3
		*boosdlex.(*boosdLex).file = $$
	}
;

imports: {}
|	imports import
	{
		$$ = append($1, $2)
	}
;

import:	YIMPORT lit ';'
	{
	}
;

kinds:	{}
|	kinds kind
	{
	}
;

kind:	YKIND id_list opt_kind ';'
	{
	}
;

opt_kind: {
		$$ = nil
	}
|	YKIND_DECL
	{
		$$ = &BasicLit{Kind:token.STRING, Value:$1.val}
	}
;

id_list: ident
	{
		$$ = []*Ident{$1}
	}
|	id_list ',' ident
	{
		$$ = append($1, $3)
	}
;

defs:	{}
|	defs def
	{
		$$ = append($1, $2)
	}
;

def:	ident top_type opt_kind specializes '{' stmts '}' ';'
	{
		if $2.val == "model" {
			$$ = &ModelDecl{Name:$1, Body:$6}
		} else {
			$$ = &InterfaceDecl{Name:$1, Body:$6}
		}
	}
;

top_type: YMODEL
	{
		$$ = $1
	}
|	YINTERFACE
	{
		$$ = $1
	}
;

specializes: {}
|	YSPECIALIZES ident
	{
		$$ = $2
	}
;


stmts:	{
		$$ = &BlockStmt{List:[]Stmt{}}
	}
|	stmts stmt
	{
		$$ = $1
		$$.List = append($1.List, $2)
	}
;

stmt:	var_decl ';'
	{
		$$ = &DeclStmt{$1}
	}
|	var_decl assignment ';'
	{
		$$ = &AssignStmt{Lhs:$1, Rhs:$2}
	}
;


var_decl:	ident opt_kind
	{
		$$ = &VarDecl{Name:$1, Units:$2}
	}
|	ident ident opt_kind
	{
		$$ = &VarDecl{Name:$1, Type:$2, Units:$3}
	}
;

assignment: '=' '{' initializers '}'
	{
		$$ = &CompositeLit{Type:NewIdent("stock"), Elts:$3}
	}
|	'=' ident '{' initializers '}'
	{
		$$ = &CompositeLit{Type:$2, Elts:$4}
	}
|	'=' expr_w_unit
	{
		$$ = $2
	}
|	'=' lit
	{
		$$ = $2
	}
;

initializers: {
		$$ = []Expr{}
	}
|	initializers initializer
	{
		$$ = append($1, $2)
	}
;

initializer: ident ':' expr_w_unit ';'
	{
		$$ = &KeyValueExpr{Key:$1, Value:$3}
	}
;

expr_w_unit: expr opt_kind
	{
		$$ = &UnitExpr{$1, $2}
	}
;

expr:	'(' expr ')'
	{
		$$  =  $2
	}
|	expr '+' expr
	{
		$$ = &BinaryExpr{X:$1, Y:$3, Op:token.ADD}
	}
|	expr '-' expr
	{
		$$ = &BinaryExpr{X:$1, Y:$3, Op:token.SUB}
	}
|	expr '*' expr
	{
		$$ = &BinaryExpr{X:$1, Y:$3, Op:token.MUL}
	}
|	expr '/' expr
	{
		$$ = &BinaryExpr{X:$1, Y:$3, Op:token.QUO}
	}
|	expr '^' expr
	{
		$$ = &BinaryExpr{X:$1, Y:$3, Op:token.XOR}
	}
|	'-' expr %prec UMINUS
	{
		$$ = &UnaryExpr{X:$2, Op:token.SUB}
	}
|	ident '(' expr_list ')' %prec FN_CALL
	{
		$$ = &CallExpr{Fun:$1, Args:$3}
	}
|	table '[' expr ']' %prec FN_CALL
	{
		$$ = &IndexExpr{X:$1, Index:$3}
	}
|	ident '[' expr ']' %prec FN_CALL
	{
		$$ = &IndexExpr{X:$1, Index:$3}
	}
|	table
	{
		$$ = $1
	}
|	ref
	{
		$$ = $1
	}
|	number
	{
		$$ = $1
	}
;

ref: ident
	{
		$$ = &RefExpr{*$1}
	}

ident:	YIDENT
	{
		$$ = &Ident{Name:$1.val}
	}
;

lit:	YLITERAL
	{
		$$ = &BasicLit{Kind:token.STRING, Value:$1.val}
	}
;

number:	YNUMBER
	{
		$$ = &BasicLit{Kind:token.FLOAT, Value:$1.val}
	}
;

expr_list: expr
	{
		$$ = make([]Expr, 1, 16)
		$$[0] = $1
	}
|	expr_list ',' expr
	{
		$$ = append($1, $3)
	}
;

table:	'[' pairs ']'
	{
		$$ = &TableExpr{Pairs: $2}
	}
;

pairs:	pair
	{
		$$ = make([]Expr, 1, 16)
		$$[0] = $1
	}
|	pairs ',' pair
	{
		$$ = append($1, $3)
	}
;

pair:	'(' number ',' number ')'
	{
		$$ = &PairExpr{$2, $4}
	}
;

%% /* start of programs */

func Parse(f *token.File, str string) *File {
	// this is weird, but without passing in a reference to this
	// result object, there isn't another good way to keep the
	// parser and lexer reentrant.
	result := &File{}
	err := boosdParse(newBoosdLex(str, f, result))
	if err != 0 {
		return nil
	}

	return result
}
