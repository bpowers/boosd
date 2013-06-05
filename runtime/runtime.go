// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"log"
	"sort"
)

type chanReq struct {
	sim    Sim
	name   string
	result chan<- float64
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
	return <-result
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
	Var(name string) (Var, bool)
	VarInfo(name string) map[string]interface{}
	Default(name string) (float64, bool)
}

// Init initializes the boosd runtime.
func Main(m Model) {
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

	for _, v := range sim.Model().VarNames() {
		vv, ok := sim.Model().Var(v)
		if !ok {
			log.Fatalf("sim.Model().Var(%s): not ok", v)
		}
		if vv.Type == TyTable {
			continue
		}

		qualName := fmt.Sprintf("%s.%s", "main", v)
		data, err := sim.ValueSeries(v)
		if err != nil {
			log.Fatalf("sim.ValueSeries(%s): %s", v, err)
		}
		series[qualName] = data[1]
		orderedVars = append(orderedVars, qualName)
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
