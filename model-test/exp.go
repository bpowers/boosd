// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// this is work on figuring out what we need to codegen
package main

import (
	"github.com/bpowers/boosd/runtime"
)

var (
	mdlMainName = "main"
	mdlMainVars = map[string]runtime.Var{

		"accum": runtime.Var{"accum", runtime.TyStock},
		"in":    runtime.Var{"in", runtime.TyFlow},
		"rate":  runtime.Var{"rate", runtime.TyAux},
	}
)

type simMain struct {
	runtime.BaseSim
}

type mdlMain struct {
	runtime.BaseModel
}

func simMainStep(s *runtime.BaseSim, dt float64) {

	s.Curr["in"] = ((s.Curr["rate"]) * (s.Curr["accum"]))

	s.Next["rate"] = s.Curr["rate"]
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
	s.Step = simMainStep

	s.Curr["rate"] = 0.070000
	s.Curr["accum"] = 200.000000

	s.Curr["time"] = ts.Start

	runtime.RegisterSim(mdlMainName, s)

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
