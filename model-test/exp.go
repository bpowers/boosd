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
	s := new(simMain)
	s.Parent = m
	return s
}

func init() {
	runtime.Register(new(mdlMain))
}

func main() {
	runtime.Main()
}
