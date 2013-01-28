// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boosd

import (
	"bytes"
	"fmt"
	"github.com/bpowers/boosd/runtime"
	"go/ast"
	"go/parser"
	"go/token"
	"text/template"
)

const modelTmpl = `
package main

import (
	"github.com/bpowers/boosd/runtime"
)

var (
	mdlMainName = "main"
	mdlMainVars = map[string]runtime.Var{
{{range $.Vars}}
		"{{.Name}}": runtime.Var{"{{.Name}}", runtime.TyVar},
{{end}}
	}
)

type simMain struct {
	runtime.BaseSim
}

type mdlMain struct {
	runtime.BaseModel
}

func simMainStep(s *runtime.BaseSim, dt float64) {
{{range $.Equations}}
	{{.}}
{{end}}
}

func (m *mdlMain) NewSim() runtime.Sim {
	ts := runtime.Timespec{
		Start:    {{$.Time.Start}},
		End:      {{$.Time.End}},
		DT:       {{$.Time.DT}},
		SaveStep: {{$.Time.SaveStep}},
	}
	tables := map[string]runtime.Table{}
	consts := runtime.Data{}

	s := new(simMain)
	s.Init(m, ts, tables, consts)
	s.Step = simMainStep

	// Initialize any constant expressions, stock initials, or
	// variables

{{range $.Initials}}
	{{.}}
{{end}}
	s.Curr["time"] = ts.Start

	runtime.RegisterSim("main", s)

	return s
}

func init() {
	m := &mdlMain{
		runtime.BaseModel{
			MName: mdlMainName,
			Vars:  mdlMainVars,
		},
	}

	runtime.RegisterModel(m)
}

func main() {
	runtime.Main()
}
`

type generator struct {
	Vars      map[string]runtime.Var
	Time      runtime.Timespec
	Equations []string
	Initials  []string
}

func (g *generator) declList(list []Decl) {
}

func (g *generator) file(f *File) []byte {
	var buf bytes.Buffer
	tmpl := template.New("model.go")
	if _, err := tmpl.Parse(modelTmpl); err != nil {
		panic(fmt.Sprintf("Parse(modelTmpl): %s", err))
	}
	if err := tmpl.Execute(&buf, g); err != nil {
		panic(fmt.Sprintf("Execute(%v): %s", g, err))
	}

	return buf.Bytes()
}

func GenGo(f *File) (*ast.File, error) {
	g := &generator{}
	code := g.file(f)

	fset := token.NewFileSet()
	goFile, err := parser.ParseFile(fset, "model.go", code, 0)
	if err != nil {
		return nil, err
	}

	return goFile, nil
}
