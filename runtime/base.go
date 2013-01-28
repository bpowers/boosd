// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	//"fmt"
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

type BaseSim struct {
	Parent Model
	Time   Timespec

	Series []Data

	Tables    map[string]Table
	Curr      Data
	Next      Data
	Constants Data

	saveEvery int64
	stepNum   int64

	// Calc{Flows,Stocks} are used by the RunTo* functions
	Step func(s *BaseSim, dt float64)
}

func (s *BaseSim) Init(m Model, ts Timespec, tables map[string]Table, consts Data) {
	s.Parent = m
	s.Time = ts

	capSeries := int((ts.End-ts.Start)/ts.SaveStep) + 1
	s.Series = make([]Data, 0, capSeries)

	s.Tables = tables
	s.Constants = consts

	// round to the nearest integer
	s.saveEvery = int64(ts.SaveStep/ts.DT + .5)

	s.Curr = Data{}
	s.Next = Data{}

	// initialize Curr with constants
	for k, v := range s.Constants {
		s.Curr[k] = v
	}
}

func (s *BaseSim) Model() Model {
	return s.Parent
}

// RunTo currently implements the Euler method
func (s *BaseSim) RunTo(t float64) error {
	for s.Curr["time"] <= t {
		s.Step(s, s.Time.DT)

		if s.stepNum % s.saveEvery == 0 {
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
	return s.RunTo(s.Time.End)
}

func (s *BaseSim) GetValue(name string) float64 {
	return -1
}

func (s *BaseSim) GetValueSeries(name string) [][2]float64 {
	return nil
}

func (s *BaseSim) SetValue(name string, val float64) error {
	return nil
}

// Step is the internal function that does one round of whatever
// integration method is in use.
//
// TODO: implement more than just euler.
func (s *BaseSim) step() {

}

type BaseModel struct{
	MName string
}

func (m *BaseModel) Name() string {
	return m.MName
}

func (m *BaseModel) Attr(name string) interface{} {
	return nil
}

func (m *BaseModel) VarNames() []string {
	return nil
}

func (m *BaseModel) VarInfo(name string) map[string]interface{} {
	return nil
}
