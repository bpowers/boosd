// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
)

type Timespec struct {
	Start    float64
	End      float64
	DT       float64
	SaveStep float64
}

// TODO: define useful methods on table
type Table [2][]float64
type Data map[string]float64
type ModelMap map[string]map[string]string
type VarMap map[string]Var
type DefaultMap map[string]float64

func (t Table) Lookup(index float64) float64 {
	size := len(t[0])
	if size == 0 {
		return 0
	}

	x := t[0]
	y := t[1]

	if index < x[0] {
		return y[0]
	} else if index > x[size-1] {
		return y[size-1]
	}

	// binary search
	low, mid, high := 0, 0, size
	for low < high {
		mid = low + (high-low)/2
		if x[mid] < index {
			low = mid + 1
		} else {
			high = mid
		}
	}

	// at this point low == high, so use 'i' for readability
	i := low
	if x[i] == index {
		return y[i]
	}

	slope := (y[i] - y[i-1]) / (x[i] - x[i-1])
	return (index-x[i-1])*slope + y[i-1]
}

type BaseSim struct {
	Parent Model
	Time   Timespec

	Coord Coordinator

	InstanceName string
	VarNames     map[string]string
	SubSims      map[string]BaseSim

	Series []Data

	Tables map[string]Table
	Curr   Data
	Next   Data

	timeSeries []float64

	saveEvery int64
	stepNum   int64

	// Calc{Flows,Stocks} are used by the RunTo* functions
	CalcInitial func(dt float64)
	CalcFlows   func(dt float64)
	CalcStocks  func(dt float64)
}

// max returns the int64 max of two numbers
func max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func (s *BaseSim) Init(m Model, ts Timespec, tables map[string]Table) {
	s.Parent = m
	s.Time = ts

	if ts.SaveStep == 0 {
		panic("can't have 0 SaveStep")
	}

	capSeries := int((ts.End-ts.Start)/ts.SaveStep) + 1
	s.Series = make([]Data, 0, capSeries)
	s.timeSeries = make([]float64, 0, capSeries)

	s.Tables = tables

	// round to the nearest integer, but make sure we're non-zero
	s.saveEvery = max(int64(ts.SaveStep/ts.DT + .5), 1)

	s.Curr = Data{}
	s.Next = Data{}

	s.Curr["time"] = ts.Start
}

func (s *BaseSim) Model() Model {
	return s.Parent
}

// RunTo currently implements the Euler method
func (s *BaseSim) RunTo(t float64) error {
	if s.Curr["time"] == s.Time.Start {
		s.CalcInitial(s.Time.DT)
	}

	for s.Curr["time"] <= t {
		s.CalcFlows(s.Time.DT)
		s.CalcStocks(s.Time.DT)

		if s.stepNum%s.saveEvery == 0 {
			s.timeSeries = append(s.timeSeries, s.Curr["time"])
			s.Series = append(s.Series, s.Curr)
		}
		s.stepNum++

		s.Next["time"] = s.Curr["time"] + s.Time.DT
		s.Curr = s.Next
		s.Next = Data{}
	}
	return nil
}

func (s *BaseSim) RunToEnd() error {
	return s.RunTo(s.Time.End + .5*s.Time.DT)
}

func (s *BaseSim) Value(name string) (v float64, err error) {
	err = fmt.Errorf("unknown var %s", name)
	return
}

func (s *BaseSim) ValueSeries(name string) (r [2][]float64, err error) {
	r[0] = s.timeSeries
	r[1] = make([]float64, len(s.Series))
	for i, d := range s.Series {
		r[1][i] = d[name]
	}
	return
}

func (s *BaseSim) SetValue(name string, val float64) error {
	return nil
}

type BaseModel struct {
	MName    string
	Vars     VarMap
	Defaults DefaultMap
	Tables   map[string]Table
}

func (m *BaseModel) Default(name string) (v float64, ok bool) {
	v, ok = m.Defaults[name]
	return
}

func (m *BaseModel) Name() string {
	return m.MName
}

func (m *BaseModel) Attr(name string) interface{} {
	return nil
}

func (m *BaseModel) Var(name string) (Var, bool) {
	v, ok := m.Vars[name]
	return v, ok
}

func (m *BaseModel) VarNames() []string {
	names := make([]string, 0, len(m.Vars))
	for n, _ := range m.Vars {
		names = append(names, n)
	}
	return names
}

func (m *BaseModel) VarInfo(name string) map[string]interface{} {
	return nil
}

type VarType int

const (
	TyAux     VarType = iota
	TyStock   VarType = iota
	TyFlow    VarType = iota
	TyTable   VarType = iota
	TyConst   VarType = iota
	TyUnknown VarType = iota
)

type Var struct {
	Name string
	Type VarType
}

var tyPretty = map[VarType]string{
	TyAux:   "TyAux",
	TyStock: "TyStock",
	TyFlow:  "TyFlow",
	TyTable: "TyTable",
	TyConst: "TyConst",
}

func (vt VarType) String() string {
	s, ok := tyPretty[vt]
	if !ok {
		s = "TyUnknown"
	}
	return s
}

var tyNames = map[string]VarType{
	"aux":   TyAux,
	"stock": TyStock,
	"flow":  TyFlow,
	"table": TyTable,
	"const": TyConst,
}

func TypeForName(n string) VarType {
	ty, ok := tyNames[n]
	if !ok {
		fmt.Printf("unknown runtime type '%s'\n", n)
		ty = TyUnknown
	}
	return ty
}

func (vm VarMap) InstanceMap(iName string) map[string]string {
	mm := map[string]string{}
	for k, _ := range vm {
		mm[k] = fmt.Sprintf("%s.%s", iName, k)
	}
	return mm
}
