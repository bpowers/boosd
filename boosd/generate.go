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
	"unicode"
)

const fileTmpl = `
{{define "modelTmpl"}}
var (
	mdl{{$.CamelName}}Name = "{{$.Name}}"
	mdl{{$.CamelName}}Vars = map[string]runtime.Var{
{{range $.Vars}}
		"{{.Name}}": runtime.Var{"{{.Name}}", runtime.{{.Type}}},{{end}}
	}
)

type sim{{$.CamelName}} struct {
	runtime.BaseSim
}

type mdl{{$.CamelName}} struct {
	runtime.BaseModel
}

func sim{{$.CamelName}}Step(s *runtime.BaseSim, dt float64) {
{{range $.Equations}}
	{{.}}{{end}}
{{range $.Stocks}}
	{{.}}{{end}}
}

func (m *mdl{{$.CamelName}}) NewSim() runtime.Sim {
	ts := runtime.Timespec{
		Start:    {{$.Time.Start}},
		End:      {{$.Time.End}},
		DT:       {{$.Time.DT}},
		SaveStep: {{$.Time.SaveStep}},
	}
	tables := map[string]runtime.Table{}
	consts := runtime.Data{}

	s := new(sim{{$.CamelName}})
	s.Init(m, ts, tables, consts)
	s.Step = sim{{$.CamelName}}Step
{{range $.Initials}}
	{{.}}{{end}}

	s.Curr["time"] = ts.Start

	runtime.RegisterSim(mdl{{$.CamelName}}Name, s)

	return s
}
{{end}}

package main

import (
	"github.com/bpowers/boosd/runtime"
)

{{range $.Models}}{{template "modelTmpl" .}}{{end}}

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

type genModel struct {
	Name      string
	CamelName string // camelcased
	Vars      map[string]runtime.Var
	Time      runtime.Timespec
	Equations []string
	Stocks    []string
	Initials  []string
	Abstract  bool
}

type generator struct {
	Models map[string]*genModel
	curr *genModel
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
			g.curr.Time.Start = v
		case "end":
			g.curr.Time.End = v
		case "dt":
			g.curr.Time.DT = v
		case "save_step":
			g.curr.Time.SaveStep = v
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

func (g *generator) initial(name string, expr Expr) (err error) {
	val, err := constEval(expr)
	if err == nil {
		init := fmt.Sprintf(`s.Curr["%s"] = %f`, name, val)
		g.curr.Initials = append(g.curr.Initials, init)
	} else {
		unitExpr, ok := expr.(*UnitExpr)
		if !ok {
			return fmt.Errorf("initial(%s): not unitexpr", name)
		}
		switch e := unitExpr.X.(type) {
		case *RefExpr:
			if _, ok := g.curr.Vars[e.Ident.Name]; ok {
				init := fmt.Sprintf(`s.Curr["%s"] = s.Curr["%s"]`,
					name, e.Name)
				g.curr.Initials = append(g.curr.Initials, init)
				err = nil
			} else {
				err = fmt.Errorf("initial(%s): non-const ident %v",
					name, e.Name)
			}
		default:
			err = fmt.Errorf("initial(%s): non-const %v (%T)",
				name, unitExpr.X, unitExpr.X)
		}
	}
	return
}

func (g *generator) stock(name string, expr Expr) error {
	cl, ok := expr.(*CompositeLit)
	if !ok {
		panic(fmt.Sprintf("stock is %T, not CompositeLit", expr))
	}
	var bi, in, out string
	for _, e := range cl.Elts {
		k, val, err := kvConvert(e)
		if err != nil {
			panic(fmt.Sprintf("kvConvert(%s): %s", name, err))
		}
		switch k {
		case "initial":
			if err := g.initial(name, val); err != nil {
				return fmt.Errorf("initial(%s, %s): %s",
					name, val, err)
			}
		case "biflow":
			bi = fmt.Sprintf("+%s", val)
		case "inflow":
			in = fmt.Sprintf("+%s", val)
		case "outflow":
			out = fmt.Sprintf("-(%s)", val)
		default:
			panic(fmt.Sprintf("stock(%s): unknown k %s",
				name, k))
		}
	}
	eqn := fmt.Sprintf(`s.Next["%s"] = s.Curr["%s"] + (%s %s %s)*dt`, name, name, bi, in, out)
	g.curr.Stocks = append(g.curr.Stocks, eqn)
	return nil
}

func (g *generator) expr(name string, expr Expr) {
	if g.curr.Vars[name].Type == runtime.TyAux {
		if isConst(expr) {
			g.initial(name, expr)
			eqn := fmt.Sprintf(`s.Next["%s"] = s.Curr["%s"]`, name, name)
			g.curr.Stocks = append(g.curr.Stocks, eqn)
		} else {
			eqn := fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
			g.curr.Equations = append(g.curr.Equations, eqn)
		}
	} else {
		eqn := fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
		g.curr.Equations = append(g.curr.Equations, eqn)
	}
}

func (g *generator) assign(s *AssignStmt) error {
	if s.Lhs.Name.Name == "timespec" {
		c, ok := s.Rhs.(*CompositeLit)
		if !ok {
			return fmt.Errorf("timespec is %T, not CompositeLit",
				s.Rhs)
		}
		g.timespec(c.Elts)
		return nil
	}
	v, ok := g.curr.Vars[s.Lhs.Name.Name]
	if !ok {
		return fmt.Errorf("assign: unknown v '%s'?", s.Lhs.Name.Name)
	}
	if v.Type == runtime.TyStock {
		if err := g.stock(v.Name, s.Rhs); err != nil {
			return err
		}
	} else {
		g.expr(v.Name, s.Rhs)
	}
	return nil
}

func (g *generator) stmt(s Stmt) error {
	switch ss := s.(type) {
	case *AssignStmt:
		if err := g.assign(ss); err != nil {
			return err
		}
	case *DeclStmt:
	default:
		log.Printf("stmt(%T): unimplemented - %T", s, s)
	}
	return nil
}

func (g *generator) vars(stmts ...Stmt) (err error) {
	addVar := func(vd *VarDecl) error {
		v, err := varFromDecl(vd)
		if err != nil {
			return fmt.Errorf("varFromDecl(%v): %s", vd, err)
		}
		g.curr.Vars[v.Name] = v
		return nil
	}
	for _, s := range stmts {
		switch ss := s.(type) {
		case *AssignStmt:
			err = addVar(ss.Lhs)
		case *DeclStmt:
			g.curr.Abstract = true
			err = addVar(ss.Decl)
		default:
			err = fmt.Errorf("varFromDecl(%v): unknown ty %T",
				s, ss)
		}
	}
	return
}

func (g *generator) model(m *ModelDecl) error {
	name := m.Name.Name
	camelName := fmt.Sprintf("%c%s", unicode.ToUpper(rune(name[0])), name[1:])
	g.curr = &genModel{
		Name: name,
		CamelName: camelName,
		Vars: map[string]runtime.Var{},
		Equations: []string{},
		Stocks: []string{},
		Initials: []string{},
	}
	g.vars(m.Body.List...)
	for _, s := range m.Body.List {
		if err := g.stmt(s); err != nil {
			return err
		}
	}
	g.Models[m.Name.Name] = g.curr
	g.curr = nil

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
	if _, err := tmpl.Parse(fileTmpl); err != nil {
		panic(fmt.Sprintf("Parse(modelTmpl): %s", err))
	}
	if err := tmpl.Execute(&buf, g); err != nil {
		panic(fmt.Sprintf("Execute(%v): %s", g, err))
	}

	return buf.Bytes(), nil
}

func GenGo(f *File) (*ast.File, error) {
	g := &generator{
		Models: map[string]*genModel{},
	}

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
