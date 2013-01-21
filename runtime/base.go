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

type BaseSim struct {
	Parent Model
	Time   Timespec

	Series []Data

	Tables    map[string]Table
	Curr      Data
	Next      Data
	Constants Data

	// Calc{Flows,Stocks} are used by the RunTo* functions
	CalcFlows  func(s BaseSim, dt float64) error
	CalcStocks func(s BaseSim, dt float64) error
}

func (s *BaseSim) Init(m Model, ts Timespec, tables map[string]Table, consts Data) {
	s.Parent = m
	s.Time = ts

	capSeries := int((ts.End-ts.Start)/ts.SaveStep) + 1
	s.Series = make([]Data, 0, capSeries)

	s.Tables = tables
	s.Constants = consts

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

func (s *BaseSim) RunTo(t float64) error {
	return fmt.Errorf("not implemented")
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
func (s *BaseSim) Step() {

}

type BaseModel struct{}

func (m *BaseModel) Attr(name string) interface{} {
	return nil
}

func (m *BaseModel) VarNames() []string {
	return nil
}

func (m *BaseModel) VarInfo(name string) map[string]interface{} {
	return nil
}
