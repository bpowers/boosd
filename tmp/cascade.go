package main

import (
	"fmt"
)

/*
model CascadedLevels specializes Tub {
	flow inflow
	flow toLevel2 = level1 * rate
	flow toLevel3 = level2 * rate
	flow outflow  = level3 * rate

	stock level1 = {
		inflow:  inflow
		outflow: toLevel2
	}
	stock level2 = {
		inflow:  toLevel2
		outflow: toLevel3
	}
	stock level3 = {
		inflow:  toLevel2
		outflow: outflow
	}

	// aggregated reporting variable
	level = level1 + level2 + level3
}
*/

type Smooth3I struct {
	input, initial, delay *float64
	ChangeIn1 float64
	ChangeIn2 float64
	ChangeIn3 float64
	Level1   float64
	Level2   float64
	Smooth3I float64
}

func (s *Smooth3I) Step(dt float64) (smooth3i float64) {

	s.ChangeIn1 = (*s.input - s.Level1) / *s.delay
	s.ChangeIn2 = (s.Level1 - s.Level2) / *s.delay
	s.ChangeIn3 = (s.Level2 - s.Smooth3I) / *s.delay

	smooth3i = s.Smooth3I

	s.Level1 += (s.ChangeIn1)*dt
	s.Level2 += (s.ChangeIn2)*dt
	s.Smooth3I += (s.ChangeIn3)*dt

	return
}

func NewSmooth3I(input, initial, delay *float64) *Smooth3I {
	result := Smooth3I{
		input: input,
		initial: initial,
		delay: delay,
	}

	// initialize levels
	result.Level1 = *initial
	result.Level2 = *initial
	result.Smooth3I = *initial

	return &result
}

func main() {
	dt := 1.0
	initialLevel := 0.0
	level1 := initialLevel
	level2 := initialLevel
	level3 := initialLevel

	inflow := 10.0
	rate := 10.0 / 3

	input := 5.0
	initial := 0.0
	delay := 2.0

	infoM := NewSmooth3I(&input, &initial, &delay)

	for i := 0.0; i <= 100; i += dt {
		info := infoM.Step(dt)

		toLevel2 := level1 / rate
		toLevel3 := level2 / rate
		outflow := level3 / rate

		fmt.Printf("%.0f\t%.2f\n", i, info)

		level1 += (inflow - toLevel2) * dt
		level2 += (toLevel2 - toLevel3) * dt
		level3 += (toLevel3 - outflow) * dt
	}
}
