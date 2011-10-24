// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// based off of Appendix A from http://dinosaur.compilertools.net/yacc/

%{

package main

import (
	"bufio"
	"fmt"
	"os"
)

var regs = map[string]int{}
var base int

%}



// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	kind int
	val  int
	id   string
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> expr

// same for terminals
%token <val> NUMBER
%token <id> IDENTIFIER

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */

%%

list	:
	|    list stat '\n'
	;


stat	:    expr
		{
			fmt.Printf( "%d\n", $1 );
		}
	|    IDENTIFIER '=' expr
		{
			regs[$1]  =  $3
		}
	;

expr	:    '(' expr ')'
		{ $$  =  $2 }
	|    expr '+' expr
		{ $$  =  $1 + $3 }
	|    expr '-' expr
		{ $$  =  $1 - $3 }
	|    expr '*' expr
		{ $$  =  $1 * $3 }
	|    expr '/' expr
		{ $$  =  $1 / $3 }
	|    expr '%' expr
		{ $$  =  $1 % $3 }
	|    expr '&' expr
		{ $$  =  $1 & $3 }
	|    expr '|' expr
		{ $$  =  $1 | $3 }
	|    '-'  expr        %prec  UMINUS
		{ $$  = -$2  }
	|    IDENTIFIER
		{ $$  = regs[$1] }
	|    NUMBER
	;

%%      /*  start  of  programs  */

func Parse(eqn string) int {
	return calcParse(newCalcLex(eqn))
}

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var eqn string
		var ok bool

		fmt.Printf("equation: ")
		if eqn, ok = readline(fi); ok {
			Parse(eqn)
		} else {
			break
		}
	}
}

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}
