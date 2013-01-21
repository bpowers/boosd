package runtime

import (
	"fmt"
)

type BaseSim struct {
	Parent Model
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
