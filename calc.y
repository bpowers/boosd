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

%}


// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	kinds []kind
	kind kind
	tok  tok
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <kind> kind
%type <kinds> kinds

// same for terminals
%token <tok> KIND ID KIND_DECL NUMBER

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */

%%

kinds:
	kind
	{
		fmt.Printf("%#v\n", $1);
//		$$ = append($1, $2);
		$$ = []kind{$1}
	}
;

kind:
	KIND ID KIND_DECL '\n'
	{
		$$  =  kind{$2.val, $3.val}
	}
//|	KIND ID '\n'
;

%% /* start of programs */

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
