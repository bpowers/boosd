package main

import (
	"boosd/runtime"
)

type simMain struct {
	runtime.BaseSim
}

type mdlMain struct {
	runtime.BaseModel
}

func (m *mdlMain) NewSim() runtime.Sim {
	ts := runtime.Timespec{
		Start:    0,
		End:      50,
		DT:       .5,
		SaveStep: 1,
	}
	tables := map[string]runtime.Table{}
	consts := runtime.Data{}

	s := new(simMain)
	s.Init(m, ts, tables, consts)

	// Initialize any constant expressions, stock initials, or
	// variables

	s.Curr["accum"] = 200
	s.Curr["rate"] = .07

	return s
}

func init() {
	m := new(mdlMain)

	runtime.Register(m)
}

func main() {
	runtime.Main()
}
