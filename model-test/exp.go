// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// this is work on figuring out what we need to codegen
package main

import (
	"github.com/bpowers/boosd/runtime"
)

var (
	mMain = mdlMain{
		runtime.BaseModel{
			MName: "main",
			Vars: runtime.VarMap{
				"accum": runtime.Var{"accum", runtime.TyStock},
				"in":    runtime.Var{"in", runtime.TyFlow},
				"rate":  runtime.Var{"rate", runtime.TyAux},
			},
			Defaults: runtime.DefaultMap{
				"accum": 200.000000,
				"rate":  0.070000,
			},
		},
	}
)

type simMain struct {
	runtime.BaseSim
}

type mdlMain struct {
	runtime.BaseModel
}

func (s *simMain) calcInitial(c runtime.Coordinator, dt float64) {
	s.Curr["accum"] = c.Data(s, "accum")
	s.Curr["rate"] = c.Data(s, "rate")
}

func (s *simMain) calcFlows(c runtime.Coordinator, dt float64) {
	s.Curr["rate"] = c.Data(s, "rate")
	s.Curr["in"] = ((s.Curr["rate"]) * (s.Curr["accum"]))
}

func (s *simMain) calcStocks(c runtime.Coordinator, dt float64) {
	s.Next["accum"] = s.Curr["accum"] + (+s.Curr["in"])*dt
}

func (m *mdlMain) NewSim(name string) runtime.Sim {
	ts := runtime.Timespec{
		Start:    0,
		End:      50,
		DT:       0.1,
		SaveStep: 1,
	}
	tables := map[string]runtime.Table{}
	consts := runtime.Data{}

	s := new(simMain)
	s.InstanceName = name
	s.Init(m, ts, tables, consts)

	s.CalcInitial = s.calcInitial
	s.CalcFlows = s.calcFlows
	s.CalcStocks = s.calcStocks

	return s
}

func main() {
	runtime.Main(&mMain)
}
