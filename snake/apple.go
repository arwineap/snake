package snake

import (
	"time"
)

type Apple struct {
	Positions     []Point
	Width         int
	LastDrop      time.Time
	DropFrequency time.Duration
}

func (a *Apple) CheckDrop() bool {
	if time.Since(a.LastDrop) < a.DropFrequency {
		return false
	}
	return true
}

func (a *Apple) Drop(position Point) {
	a.Positions = append(a.Positions, position)
	a.LastDrop = time.Now()
}

func (a *Apple) Remove(p Point) {
	var removeIdx int
	var exists bool
	for i, ap := range a.Positions {
		if p == ap {
			removeIdx = i
			exists = true
			break
		}
		exists = false
	}
	if exists {
		a.Positions[removeIdx] = a.Positions[len(a.Positions)-1]
		a.Positions = a.Positions[:len(a.Positions)-1]
	}
}
