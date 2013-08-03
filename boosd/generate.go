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
	"strings"
	"text/template"
	"unicode"
)

const fileTmpl = `
{{define "modelTmpl"}}
var m{{$.CamelName}} = mdl{{$.CamelName}}{
	runtime.BaseModel{
		MName: "{{$.Name}}",
		Vars: runtime.VarMap{ {{range $.Vars}}
			"{{.Name}}": runtime.Var{"{{.Name}}", runtime.{{.Type}}},{{end}}
		},
		Defaults: runtime.DefaultMap{ {{range $n, $_ := $.Initials}}
			{{if simple . }}"{{$n}}": {{.}}, {{end}}{{end}}
		},
		Tables: map[string]runtime.Table{ {{range $n, $_ := $.Tables}}
			"{{$n}}": {{printf "%#v" .}}, {{end}}
		},
	},
}

type sim{{$.CamelName}} struct {
	runtime.BaseSim
}

type mdl{{$.CamelName}} struct {
	runtime.BaseModel
}

func (s *sim{{$.CamelName}}) calcInitial(dt float64) { {{if $.Initials }}
	c := s.Coord
	{{end}} {{range $n, $_ := $.Initials}}
	s.Curr["{{$n}}"] = {{if simple .}}c.Data(s, "{{$n}}"){{else}}{{.}}{{end}}{{end}}
}

func (s *sim{{$.CamelName}}) calcFlows(dt float64) { {{if $.UseCoordFlows }}
	c := s.Coord
	{{end}} {{range $.Equations}}
	{{.}}{{end}}
}

func (s *sim{{$.CamelName}}) calcStocks(dt float64) { {{if $.UseCoordStocks }}
	c := s.Coord
	{{end}} {{range $.Stocks}}
	{{.}}{{end}}
}

func (m *mdl{{$.CamelName}}) NewSim(name string, c runtime.Coordinator) runtime.Sim {
	ts := runtime.Timespec{
		Start:    {{$.Time.Start}},
		End:      {{$.Time.End}},
		DT:       {{$.Time.DT}},
		SaveStep: {{$.Time.SaveStep}},
	}

	s := new(sim{{$.CamelName}})
	s.InstanceName = name
	s.Parent = m
	s.Coord = c

	s.Init(m, ts, m.Tables)

	s.CalcInitial = s.calcInitial
	s.CalcFlows = s.calcFlows
	s.CalcStocks = s.calcStocks

	return s
}
{{end}}

package main

import (
	"github.com/bpowers/boosd/runtime"
)

{{range $.Models}}{{template "modelTmpl" .}}{{end}}

func main() {
	runtime.Main(&mMain)
}
`

type genModel struct {
	Name      string
	CamelName string // camelcased
	Vars      map[string]runtime.Var
	Tables    map[string]runtime.Table
	Time      runtime.Timespec
	Equations []string
	Stocks    []string
	Initials  map[string]string
	Abstract  bool
	UseCoordFlows bool
	UseCoordStocks bool
}

type generator struct {
	Models map[string]*genModel
	curr   *genModel
}

func (g *generator) declList(list []Decl) {
}

// stripUnits returns the child of rhs if rhs is a UnitExpr, and
// returns rhs otherwise.
func stripUnits(rhs Expr) Expr {
	switch r := rhs.(type) {
	case *UnitExpr:
		rhs = r.X
	}
	return rhs
}

// constEval returns the float64 value represented by Expr, or an
// error if it can't be evaluated at compile time.
func constEval(e Expr) (v float64, err error) {
	// if we're wrapped in units, remove them.  Unit safety is a
	// separate issue.
	e = stripUnits(e)
	basic, ok := e.(*BasicLit)
	if !ok {
		err = fmt.Errorf("val %T not BasicLit", e)
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

func varFromDecl(d *VarDecl) (v runtime.Var, err error) {
	//log.Printf("var '%s': %s - %s", d.Name.Name, d.Type.Name, runtime.TypeForName(d.Type.Name))
	return runtime.Var{d.Name.Name, runtime.TypeForName(d.Type.Name)}, nil
}

func (g *generator) initial(name string, expr Expr) (err error) {
	val, err := constEval(expr)
	if err == nil {
		init := fmt.Sprintf(`%f`, val)
		g.curr.Initials[name] = init
	} else {
		expr = stripUnits(expr)
		switch e := expr.(type) {
		case *RefExpr:
			if _, ok := g.curr.Vars[e.Ident.Name]; ok {
				init := fmt.Sprintf(`s.Curr["%s"]`, e.Name)
				g.curr.Initials[name] = init
				err = nil
			} else {
				err = fmt.Errorf("initial(%s): non-const ident %v",
					name, e.Name)
			}
		default:
			err = fmt.Errorf("initial(%s): non-const %v (%T)",
				name, expr, expr)
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

func (g *generator) table(name string, e Expr) error {
	var t *TableExpr

	// if we're wrapped in units, remove them.  Unit safety is a
	// separate issue.
	e = stripUnits(e)

	switch r := e.(type) {
	case *TableExpr:
		t = r
	case *IndexExpr:
		t, _ = r.X.(*TableExpr)

		eqn := fmt.Sprintf(`s.Curr["%s"] = s.Tables["%s"].Lookup(%s)`,
			name, name, r.Index)
		g.curr.Equations = append(g.curr.Equations, eqn)

	default:
		return fmt.Errorf("table w/ non-table '%s': %#v", name, e)
	}

	l := len(t.Pairs)
	tab := [2][]float64{make([]float64, l), make([]float64, l)}

	for i, p := range t.Pairs {
		x, err := constEval(p.X)
		if err != nil {
			return fmt.Errorf("pair %d X (%s): %s", i, p.X, err)
		}
		y, err := constEval(p.Y)
		if err != nil {
			return fmt.Errorf("pair %d Y (%s): %s", i, p.Y, err)
		}
		tab[0][i] = x
		tab[1][i] = y
	}

	g.curr.Tables[name] = tab

	return nil
}

func (g *generator) expr(name string, expr Expr) {
	var eqn string
	switch g.curr.Vars[name].Type {
	case runtime.TyAux:
		if isConst(expr) {
			g.initial(name, expr)
			eqn = fmt.Sprintf(`s.Curr["%s"] = c.Data(s, "%s")`, name, name)
			g.curr.UseCoordFlows = true
		} else if e, ok := expr.(*CompositeLit); ok {
			log.Printf("%s - composit lit %T (%#v)", name, e, e)
		} else {
			eqn = fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
		}
	case runtime.TyTable:
		if err := g.table(name, expr); err != nil {
			log.Printf("table(%s): %s", name, err)
		}
	default:
		log.Printf("%s - expr2 %T (%#v)", name, expr, expr)
		eqn = fmt.Sprintf(`s.Curr["%s"] = %s`, name, expr)
	}
	if len(eqn) > 0 {
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
		v, ok := g.curr.Vars[ss.Decl.Name.Name]
		if !ok {
			return fmt.Errorf("assign: unknown v '%s'?", ss.Decl.Name.Name)
		}
		eqn := fmt.Sprintf(`s.Curr["%s"] = c.Data(s, "%s")`, v.Name, v.Name)
		g.curr.Equations = append(g.curr.Equations, eqn)
		g.curr.Initials[v.Name] = eqn
	default:
		log.Printf("stmt(%T): unimplemented - %T", s, s)
	}
	return nil
}

var (
	identAux   = Ident{Name: "aux"}
	identTable = Ident{Name: "table"}
)

// resolveType has two uses - if a decl was given an explicit type, it
// verifies this matches the type of the rhs expression.  If a decl
// doesn't have an explicit type, it figures out the implicit type
// from rhs.
func resolveType(d *VarDecl, rhs Expr) error {
	if d.Type != nil {
		// TODO: verify type matches rhs
		return nil
	}

	// if we're wrapped in units, remove them.  Unit safety is a
	// separate issue.
	rhs = stripUnits(rhs)

	switch r := rhs.(type) {
	case *TableExpr:
		d.Type = &identTable
	case *IndexExpr:
		if _, ok := r.X.(*TableExpr); ok {
			d.Type = &identTable
		} else {
			d.Type = &identAux
		}
	default:
		d.Type = &identAux
	}

	return nil
}

func (g *generator) vars(stmts ...Stmt) (err error) {
	addVar := func(vd *VarDecl) error {
		v, err := varFromDecl(vd)
		if err != nil {
			return fmt.Errorf("varFromDecl(%v): %s", vd, err)
		}
		if v.Name != "timespec" {
			g.curr.Vars[v.Name] = v
		}
		return nil
	}
outer:
	for i, s := range stmts {
		switch ss := s.(type) {
		case *AssignStmt:
			if err = resolveType(ss.Lhs, ss.Rhs); err != nil {
				break outer
			}
			err = addVar(ss.Lhs)
		case *DeclStmt:
			g.curr.Abstract = true
			err = addVar(ss.Decl)
		default:
			err = fmt.Errorf("stmt %d (%v): unknown ty %T",
				i, s, ss)
		}
		if err != nil {
			break
		}
	}
	return
}

func (g *generator) model(m *ModelDecl) error {
	name := m.Name.Name
	camelName := fmt.Sprintf("%c%s", unicode.ToUpper(rune(name[0])), name[1:])
	g.curr = &genModel{
		Name:      name,
		CamelName: camelName,
		Vars:      map[string]runtime.Var{},
		Tables:    map[string]runtime.Table{},
		Equations: []string{},
		Stocks:    []string{},
		Initials:  map[string]string{},
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

func tmplSimple(eqn string) bool {
	return !strings.HasPrefix(eqn, `s.Curr["`)
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
	tmpl = tmpl.Funcs(template.FuncMap{"simple": tmplSimple})
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
