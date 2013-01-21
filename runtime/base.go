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

type BaseSim struct {
	Parent   Model
	Timespec Timespec

	Tables map[string]Table
	Curr   map[string]float64
	Next   map[string]float64
	// TODO: eventually we should pull constants out into their
	// own map
	//Constants map[string]float64
}

func (s *BaseSim) Model() Model {
	return s.Parent
}

func (s *BaseSim) RunTo(t float64) error {
	return fmt.Errorf("not implemented")
}

func (s *BaseSim) RunToEnd() {
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

type BaseModel struct {
	Tables    map[string]Table
	Constants map[string]float64
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
