// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"log"
	"sort"
)

var models = map[string]Model{}
var sims = map[string]Sim{}

type chanReq struct{
	sim     Sim
	name    string
	result  chan<- float64
}

// a coordinator is the syncronization point between model instances
// and provider of gaming variables.
type Coordinator struct {
	req  chan chanReq
	stop chan struct{}
}

func (c *Coordinator) worker() {
outer:
	for {
		select {
		case req := <-c.req:
			val, _ := req.sim.Model().Default(req.name)
			req.result <- val
		case <-c.stop:
			break outer
		}
	}
}

func (c *Coordinator) Data(s Sim, name string) float64 {
	result := make(chan float64)
	c.req <- chanReq{s, name, result}
	return <- result
}

func (c *Coordinator) Close() {
	close(c.stop)
}

func NewCoordinator() Coordinator {
	c := Coordinator{
		req: make(chan chanReq),
	}
	go c.worker()
	return c
}

type Sim interface {
	Model() Model

	RunTo(t float64) error
	RunToEnd() error

	Value(name string) (float64, error)
	ValueSeries(name string) ([2][]float64, error)

	SetValue(name string, val float64) error
}

type Model interface {
	Name() string
	NewSim(iName string) Sim
	Attr(name string) interface{}
	VarNames() []string
	VarInfo(name string) map[string]interface{}
	Default(name string) (float64, bool)
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

	sim := m.NewSim("main")

	if err := sim.RunToEnd(); err != nil {
		log.Fatalf("sim.RunToEnd: %s", err)
	}

	tsRaw, err := sim.ValueSeries("time")
	if err != nil {
		log.Fatalf("sim.ValueSeries(time): %s", err)
	}
	timeSeries := tsRaw[1]
	series := map[string][]float64{}
	orderedVars := sort.StringSlice{}

	for simName, s := range sims {
		for _, v := range s.Model().VarNames() {
			qualName := fmt.Sprintf("%s.%s", simName, v)
			data, err := s.ValueSeries(v)
			if err != nil {
				log.Fatalf("s.ValueSeries(%s): %s", v, err)
			}
			series[qualName] = data[1]
			orderedVars = append(orderedVars, qualName)
		}
	}

	orderedVars.Sort()

	fmt.Printf("time")
	for _, v := range orderedVars {
		fmt.Printf("\t%s", v)
	}
	fmt.Printf("\n")

	for i, t := range timeSeries {
		fmt.Printf("%f", t)
		for _, v := range orderedVars {
			fmt.Printf("\t%f", series[v][i])
		}
		fmt.Printf("\n")
	}
}
