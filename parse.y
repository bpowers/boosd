// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// based off of Appendix A from http://dinosaur.compilertools.net/yacc/

%{

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	kinds  []kind
	kind   kind
	tok    tok
	ids    []string
	file     File
	str    string
	mdl    mdl
	models []mdl
	node   node
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <kind>   kind
%type <kinds>  kinds
%type <ids>    id_list imports callable
%type <file>   file
%type <str>    pkg import opt_kind specializes
%type <mdl>    def
%type <models> defs
%type <node>   expr
%type <tok>    top_type

// same for terminals
%token <tok> IMPORT KIND IDENT KIND_DECL NUMBER LITERAL PACKAGE MODEL
%token <tok> CALLABLE SPECIALIZES INTERFACE

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
			pkgName: $1,
			imports: $2,
			kinds: $3,
			models: $4,
		}
	}
;

pkg:	PACKAGE IDENT ';'
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

import:	IMPORT LITERAL ';'
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

kind:	KIND id_list opt_kind ';'
	{
		$$  =  kind{$2, $3}
	}
;

opt_kind: {
		$$ = ""
	}
|	KIND_DECL
	{
		$$ = $1.val
	}
;

id_list: IDENT
	{
		$$ = []string{$1.val}
	}
|	id_list ',' IDENT
	{
		$$ = append($1, $3.val)
	}
;

defs:	{}
|	defs def
	{
		$$ = append($1, $2)
	}
;

def:	top_type IDENT opt_kind callable specializes '{' stmts '}' ';'
	{
		$$.sig = $4
	}
;

top_type: MODEL
	{
		$$ = $1
	}
|	INTERFACE
	{
		$$ = $1
	}
;

callable: {}
|	CALLABLE '(' id_list ')'
	{
		$$ = $3
	}
;

specializes: {}
|	SPECIALIZES IDENT
	{
		$$ = $2.val
	}
;


stmts:	{}
|	stmts stmt
	{
	}
;

stmt:	IDENT opt_kind assignment ';'
	{
	}
|	IDENT IDENT opt_kind assignment ';'
	{
	}
;

assignment:
	{
	}
|	'=' '{' initializers '}'
	{
	}
|	'=' IDENT '{' initializers '}'
	{
	}
|	'=' expr_w_unit
	{
	}
|	'=' LITERAL
	{
	}
;

initializers: {}
|	initializers initializer
	{
	}
;

initializer: IDENT ':' expr_w_unit ';'
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
|	IDENT '(' expr_list ')' %prec FN_CALL
	{}
|	table '[' expr ']' %prec FN_CALL
	{}
|	IDENT '[' expr ']' %prec FN_CALL
	{}
|	table
	{}
|	IDENT
	{}
|	NUMBER
	{}
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

pair:	'(' NUMBER ',' NUMBER ')'
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

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))
	units, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}
	f := Parse(string(units))
	log.Printf("compilationUnit: %#v\n", f)
}
