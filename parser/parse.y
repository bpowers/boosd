// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// based off of Appendix A from http://dinosaur.compilertools.net/yacc/

%{

package main

import (
	"fmt"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	kinds  []kind
	kind   kind
	tok    tok
	strs   []string
	ids    []*Ident
	file   File
	id     *Ident
	str    string
	mdl    mdl
	models []mdl
	node   Node
	expr   Expr
	stmt   Stmt
	block  *BlockStmt
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <kind>   kind
%type <kinds>  kinds
%type <strs>   imports
%type <ids>    id_list callable
%type <file>   file
%type <str>    pkg import opt_kind
%type <mdl>    def
%type <models> defs
%type <node>   expr
%type <id>     ident specializes
%type <tok>    top_type lit
%type <block>  stmts
%type <stmt>   stmt

// same for terminals
%token <tok> YIMPORT YKIND YKIND_DECL YPACKAGE
%token <tok> YCALLABLE YSPECIALIZES YINTERFACE YMODEL
%token <tok> YIDENT YLITERAL YNUMBER

%left '+'  '-'
%left '*'  '/'
%left '^'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */
%left FN_CALL
%%

file:	pkg
	imports
	kinds
	defs
	{
		*boosdlex.(*boosdLex).file = File{
		}
	}
;

pkg:	YPACKAGE YIDENT ';'
	{
		$$ = $2.val
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
		$$ = $2.val
	}
;

kinds:	{}
|	kinds kind
	{
		$$ = append($1, $2)
	}
;

kind:	YKIND id_list opt_kind ';'
	{
		$$  =  kind{}
	}
;

opt_kind: {
		$$ = ""
	}
|	YKIND_DECL
	{
		$$ = $1.val
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

def:	top_type ident opt_kind callable specializes '{' stmts '}' ';'
	{
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

callable: {}
|	YCALLABLE '(' id_list ')'
	{
		$$ = $3
	}
;

specializes: {}
|	YSPECIALIZES ident
	{
		$$ = $2
	}
;


stmts:	{
		$$ = &BlockStmt{List:make([]Stmt, 2)}
	}
|	stmts stmt
	{
		$$ = $1
		$$.List = append($1.List, $2)
	}
;

stmt:	ident opt_kind assignment ';'
	{
	}
|	ident ident opt_kind assignment ';'
	{
	}
;

assignment:
	{
	}
|	'=' '{' initializers '}'
	{
	}
|	'=' ident '{' initializers '}'
	{
	}
|	'=' expr_w_unit
	{
	}
|	'=' lit
	{
	}
;

initializers: {}
|	initializers initializer
	{
	}
;

initializer: ident ':' expr_w_unit ';'
	{
	}
;

expr_w_unit: expr opt_kind
	{
	}
;

expr:	'(' expr ')'
	{
		$$  =  $2
	}
|	expr '+' expr
	{}
|	expr '-' expr
	{}
|	expr '*' expr
	{}
|	expr '/' expr
	{}
|	expr '^' expr
	{}
|	'-' expr %prec UMINUS
	{}
|	ident '(' expr_list ')' %prec FN_CALL
	{}
|	table '[' expr ']' %prec FN_CALL
	{}
|	ident '[' expr ']' %prec FN_CALL
	{}
|	table
	{}
|	ident
	{}
|	number
	{}
;

ident:	YIDENT
	{
		$$ = &Ident{Name:$1.val}
	}
;

lit:	YLITERAL
	{
	}
;

number:	YNUMBER
	{
	}
;

expr_list: expr
	{
	}
|	expr_list ',' expr
	{
	}
;

table:	'[' pairs ']'
	{
	}
;

pairs:	pair
	{
	}
|	pairs ',' pair
	{
	}
;

pair:	'(' number ',' number ')'
	{
	}
;

%% /* start of programs */

func Parse(str string) *File {
	f := &File{}
	err := boosdParse(newBoosdLex(str, f))
	if err != 0 {
		return nil
	}

	return f
}
