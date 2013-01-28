// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boosd

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
)

type generator struct {
	bytes.Buffer
}

func (g *generator) declList(list []Decl) {
}

func (g *generator) file(f *File) {
	g.WriteString("package main\n\n")
}

func GenGo(f *File) (*ast.File, error) {
	g := &generator{}
	g.file(f)

	fset := token.NewFileSet()
	goFile, err := parser.ParseFile(fset, "model.go", g, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return goFile, nil
}
