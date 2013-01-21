package main

import (
	"boosd/runtime"
	"fmt"
)

type simMain struct{}

func (s *simMain) Model() runtime.Model {
	return nil
}

func (s *simMain) RunTo(t float64) error {
	return fmt.Errorf("not implemented")
}

func (s *simMain) RunToEnd() {
}

func (s *simMain) GetValue(name string) float64 {
	return -1
}
func (s *simMain) GetValueSeries(name string) [][2]float64 {
	return nil
}

func (s *simMain) SetValue(name string, val float64) error {
	return nil
}

type mdlMain struct{}

func (m *mdlMain) NewSim() runtime.Sim {
	return &simMain{}
}

func (m *mdlMain) Attr(name string) interface{} {
	return nil
}

func (m *mdlMain) VarNames() []string {
	return nil
}

func (m *mdlMain) VarInfo(name string) map[string]interface{} {
	return nil
}

func init() {
	runtime.Register(new(mdlMain))
}

func main() {
	runtime.Main()
}
