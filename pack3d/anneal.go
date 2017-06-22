package pack3d

import (
	"fmt"
	"math"
	"math/rand"
)

type AnnealCallback func(Annealable)

type Annealable interface {
	Energy() float64
	DoMove() interface{}
	UndoMove(interface{})
	Copy() Annealable
}

func Anneal(state Annealable, maxTemp, minTemp float64, steps int, callback AnnealCallback) Annealable {
	factor := -math.Log(maxTemp / minTemp)
	state = state.Copy()
	bestState := state.Copy()
	if callback != nil {
		callback(bestState)
	}
	bestEnergy := state.Energy()
	previousEnergy := bestEnergy
	rate := steps / 1000
	for step := 0; step < steps; step++ {
		pct := float64(step) / float64(steps-1)
		temp := maxTemp * math.Exp(factor*pct)
		if step%rate == 0 {
			showProgress(step, steps, bestEnergy)
		}
		undo := state.DoMove()
		energy := state.Energy()
		change := energy - previousEnergy
		if change > 0 && math.Exp(-change/temp) < rand.Float64() {
			state.UndoMove(undo)
		} else {
			previousEnergy = energy
			if energy < bestEnergy {
				bestEnergy = energy
				bestState = state.Copy()
				if callback != nil {
					callback(bestState)
				}
			}
		}
	}
	showProgress(steps, steps, bestEnergy)
	fmt.Println()
	return bestState
}

func showProgress(i, n int, e float64) {
	pct := int(100 * float64(i) / float64(n))
	fmt.Printf("  %3d%% [", pct)
	for p := 0; p < 100; p += 3 {
		if pct > p {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %.8f    \r", e)
}
