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
	kinds []kind
	kind  kind
	tok   tok
	ids   []string
	cu    compilationUnit
	str   string
	mdl   mdl
	models []mdl
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <kind>  kind
%type <kinds> kinds
%type <ids>   ids imports callable
%type <cu>    file
%type <str>   pkg import
%type <mdl>   def
%type <models> defs

// same for terminals
%token <tok> IMPORT KIND ID KIND_DECL NUMBER LITERAL PACKAGE MODEL
%token <tok> CALLABLE

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */
%%

file:	pkg
	imports
	kinds
	defs
	{
		*boosdlex.(*boosdLex).cu = compilationUnit{
			pkgName: $1,
			imports: $2,
			kinds: $3,
			models: $4,
		}
	}
;

pkg:	PACKAGE ID ';'
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

kind:
	KIND ids KIND_DECL ';'
	{
		$$  =  kind{$2, $3.val}
	}
|	KIND ids ';'
	{
		$$  =  kind{$2, ""}
	}
;

ids:	ID
	{
		$$ = []string{$1.val}
	}
|	ids ',' ID
	{
		$$ = append($1, $3.val)
	}

defs:
{
}
|	defs def
	{
		$$ = append($1, $2)
	}
;

def:	MODEL ID callable '{' stmts '}' ';'
	{
		$$.sig = $3
	}
;

callable: {}
|	CALLABLE '(' ids ')'
	{
		$$ = $3
	}
;

stmts:	{}
;

%% /* start of programs */

func Parse(str string) *compilationUnit {
	cu := &compilationUnit{}
	err := boosdParse(newBoosdLex(str, cu))
	if err != 0 {
		return nil
	}

	return cu
}

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))
	units, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}
	cu := Parse(string(units))
	log.Printf("compilationUnit: %#v\n", cu)
}
