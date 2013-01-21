package runtime

import (
	"log"
)

var models = make([]Model, 0, 1)

type Sim interface {
	Model() Model

	RunTo(t float64) error
	RunToEnd()

	GetValue(name string) float64
	GetValueSeries(name string) [][2]float64

	SetValue(name string, val float64) error
}

type Model interface {
	NewSim() Sim
	Attr(name string) interface{}
	VarNames() []string
	VarInfo(name string) map[string]interface{}
}

func Register(ms ...Model) {
	models = append(models, ms...)
}

// Init initializes the boosd runtime.
func Main() {
	// FIXME: this should be temporary
	if len(models) != 1 {
		log.Fatalf("len(%v) != 1", models)
	}

	m := models[0]

	sim := m.NewSim()

	sim.RunToEnd()

	// TODO: print results
}
