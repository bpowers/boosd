// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"log"
)

var models = map[string]Model{}
var sims = map[string]Sim{}

type Sim interface {
	Model() Model

	RunTo(t float64) error
	RunToEnd() error

	GetValue(name string) float64
	GetValueSeries(name string) [][2]float64

	SetValue(name string, val float64) error
}

type Model interface {
	Name() string
	NewSim() Sim
	Attr(name string) interface{}
	VarNames() []string
	VarInfo(name string) map[string]interface{}
}

func RegisterModel(ms ...Model) {
	for _, m := range ms {
		models[m.Name()] = m
	}
}

func RegisterSim(name string, s Sim) {
	if existing, ok := sims[name]; ok {
		panic(fmt.Sprintf("sim %s already registered (%#v)",
			name, existing))
	}
	sims[name] = s
}

// Init initializes the boosd runtime.
func Main() {
	m, ok := models["main"]
	if !ok {
		log.Fatalf("no main model registered")
	}

	sim := m.NewSim()

	if err := sim.RunToEnd(); err != nil {
		log.Fatalf("sim.RunToEnd: %s", err)
	}

	for simName, s := range sims {
		for _, v := range s.Model().VarNames() {
			fmt.Printf("%s.%s\n", simName, v)
		}
	}
}
