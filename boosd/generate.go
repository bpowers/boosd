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
	"log"
	"strconv"
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
		"{{.Name}}": runtime.Var{"{{.Name}}", runtime.{{.Type}}},{{end}}
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
{{range $.Stocks}}
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
	Stocks    []string
	Initials  []string
}

func (g *generator) declList(list []Decl) {
}

// constEval returns the float64 value represented by Expr, or an
// error if it can't be evaluated at compile time.
func constEval(e Expr) (v float64, err error) {
	val, ok := e.(*UnitExpr)
	if !ok {
		err = fmt.Errorf("timespec val %T not UnitExpr", e)
		return
	}
	basic, ok := val.X.(*BasicLit)
	if !ok {
		err = fmt.Errorf("timespec val %T not BasicLit", val.X)
		return
	}
	return strconv.ParseFloat(basic.Value, 64)
}

func isConst(e Expr) bool {
	if _, err := constEval(e); err == nil {
		return true
	}
	return false
}

func kvConvert(e Expr) (k string, v Expr, err error) {
	kv, ok := e.(*KeyValueExpr)
	if !ok {
		err = fmt.Errorf("e %T not KVeyValueExpr: %v", e, e)
		return
	}
	ident, ok := kv.Key.(*Ident)
	if !ok {
		err = fmt.Errorf("e key %T not Ident: %v", kv.Key, kv.Key)
		return
	}
	return ident.Name, kv.Value, nil
}

func (g *generator) timespec(elts []Expr) {
	for _, e := range elts {
		k, val, err := kvConvert(e)
		if err != nil {
			panic(fmt.Sprintf("kvConvert(%v): %s", e, err))
		}
		v, err := constEval(val)
		if err != nil {
			panic(fmt.Sprintf("timespec constEval(%v): %s",
				val, err))
		}
		switch k {
		case "start":
			g.Time.Start = v
		case "end":
			g.Time.End = v
		case "dt":
			g.Time.DT = v
		case "save_step":
			g.Time.SaveStep = v
		default:
			panic(fmt.Sprintf("timespec unknown key %s", k))
		}
	}
}

var identAux = &Ident{Name: "aux"}

func varFromDecl(d *VarDecl) (v runtime.Var, err error) {
	if d.Type == nil {
		d.Type = identAux
	}
	return runtime.Var{d.Name.Name, runtime.TypeForName(d.Type.Name)}, nil
}

func (g *generator) initial(name string, expr Expr) {
	val, err := constEval(expr)
	if err != nil {
		panic(fmt.Sprintf("initial(%s): non-const %v", name, expr))
	}
	init := fmt.Sprintf(`s.Curr["%s"] = %f`, name, val)
	g.Initials = append(g.Initials, init)
}

func (g *generator) stock(name string, expr Expr) {
	cl, ok := expr.(*CompositeLit)
	if !ok {
		panic(fmt.Sprintf("stock is %T, not CompositeLit", expr))
	}
	var in, out string
	for _, e := range cl.Elts {
		k, val, err := kvConvert(e)
		if err != nil {
			panic(fmt.Sprintf("kvConvert(%s): %s", name, err))
		}
		switch k {
		case "initial":
			g.initial(name, val)
		case "inflow":
			in = fmt.Sprintf("%s", val)
		case "outflow":
			out = fmt.Sprintf("-(%s)", val)
		default:
			panic(fmt.Sprintf("stock(%s): unknown k %s",
				name, k))
		}
	}
	eqn := fmt.Sprintf(`s.Next["%s"] = s.Curr["%s"] + (%s %s)*dt`, name, name, in, out)
	g.Stocks = append(g.Stocks, eqn)
}

func (g *generator) expr(name string, expr Expr) {
	if g.Vars[name].Type == runtime.TyAux {
		if isConst(expr) {
			g.initial(name, expr)
			eqn := fmt.Sprintf(`s.Next["%s"] = %s`, name, expr)
			g.Stocks = append(g.Stocks, eqn)
		} else {
			eqn := fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
			g.Equations = append(g.Equations, eqn)
		}
	} else {
		eqn := fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
		g.Equations = append(g.Equations, eqn)
	}
}

func (g *generator) assign(s *AssignStmt) {
	if s.Lhs.Name.Name == "timespec" {
		c, ok := s.Rhs.(*CompositeLit)
		if !ok {
			panic(fmt.Sprintf("timespec is %T, not CompositeLit",
				s.Rhs))
		}
		g.timespec(c.Elts)
		return
	}
	v, err := varFromDecl(s.Lhs)
	if err != nil {
		panic(fmt.Sprintf("varFromDecl(%v): %s", s.Lhs, err))
	}
	g.Vars[v.Name] = v
	if v.Type == runtime.TyStock {
		g.stock(v.Name, s.Rhs)
	} else {
		g.expr(v.Name, s.Rhs)
	}
}

func (g *generator) stmt(s Stmt) {
	switch ss := s.(type) {
	case *AssignStmt:
		g.assign(ss)
	default:
		log.Printf("stmt(%T): unimplemented - %v", s, s)
	}
}

func (g *generator) model(m *ModelDecl) error {
	if m.Name.Name != "main" {
		return fmt.Errorf("non-main model (%s) not supported", m.Name.Name)
	}
	for _, s := range m.Body.List {
		g.stmt(s)
	}
	return nil
}

func (g *generator) file(f *File) ([]byte, error) {
	for _, d := range f.Decls {
		md, ok := d.(*ModelDecl)
		if !ok {
			log.Printf("top level decl that isn't a model: %v (%T)", d, d)
			continue
		}
		if err := g.model(md); err != nil {
			return nil, fmt.Errorf("g.model: %s", err)
		}
	}

	var buf bytes.Buffer
	tmpl := template.New("model.go")
	if _, err := tmpl.Parse(modelTmpl); err != nil {
		panic(fmt.Sprintf("Parse(modelTmpl): %s", err))
	}
	if err := tmpl.Execute(&buf, g); err != nil {
		panic(fmt.Sprintf("Execute(%v): %s", g, err))
	}

	return buf.Bytes(), nil
}

func GenGo(f *File) (*ast.File, error) {
	g := new(generator)
	g.Vars = map[string]runtime.Var{}
	g.Equations = []string{}
	g.Stocks = []string{}
	g.Initials = []string{}

	code, err := g.file(f)
	if err != nil {
		return nil, fmt.Errorf("g.file: %s", err)
	}
	log.Printf("c: %s", code)

	fset := token.NewFileSet()
	goFile, err := parser.ParseFile(fset, "model.go", code, 0)
	if err != nil {
		return nil, err
	}

	return goFile, nil
}
