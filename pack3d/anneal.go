package pack3d

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/fauxgl"
)

type AnnealCallback func(Annealable)

type Annealable interface {
	Energy() float64
	DoMove([]fauxgl.Vector, fauxgl.Vector, int) (Undo, int)
	UndoMove(Undo)
	Copy() Annealable
}

func Anneal(state Annealable, maxTemp, minTemp float64, steps int, callback AnnealCallback, singleStlSize []fauxgl.Vector, frameSize fauxgl.Vector, packItemNum int) (Annealable, int) {
	start := time.Now()
	factor := -math.Log(maxTemp / minTemp)
	state = state.Copy()
	bestState := state.Copy()
	if callback != nil {
		callback(bestState)
	}
	bestEnergy := state.Energy()
	previousEnergy := bestEnergy
	rate := steps / 200
	var cycleIndex int
	for step := 0; step < steps; step++ {
		pct := float64(step) / float64(steps-1)
		temp := maxTemp * math.Exp(factor*pct)
		// every 200 steps show progress
		if step%rate == 0 {
			showProgress(step, steps, bestEnergy, time.Since(start).Seconds())
		}
		undo, ntime := state.DoMove(singleStlSize, frameSize, packItemNum)
		cycleIndex = ntime
		if ntime >= 100{
			return bestState, ntime
		}
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
	showProgress(steps, steps, bestEnergy, time.Since(start).Seconds())
	fmt.Println()
	return bestState, cycleIndex
}

/* This function shows progress, not necessary*/
func showProgress(i, n int, e, d float64) {
	pct := int(100 * float64(i) / float64(n))
	fmt.Printf("  %3d%% [", pct)
	for p := 0; p < 100; p += 3 {
		if pct > p {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %.6f %.3fs    \r", e, d)
}
