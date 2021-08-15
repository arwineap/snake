package snake

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Apple struct {
	Positions     []Point
	LastDrop      time.Time
	DropFrequency time.Duration
	Width         int
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

func (a *Apple) Draw(screen *ebiten.Image) {
	for _, p := range a.Positions {
		drawPoint(screen, p, color.RGBA{R: 255, G: 100, B: 100, A: 255})
	}
}

// PointCollides - check if any apple collides with point
func (a *Apple) PointCollides(point Point) bool {
	appleCollisions := a.Collisions(point)
	if len(appleCollisions) == 0 {
		return false
	}
	return true
}

// Collisions - returns all points that collide
func (a *Apple) Collisions(point Point) []Point {
	var results []Point
	for _, p := range a.Positions {
		if checkPointCollision(p, point) {
			results = append(results, p)
		}
	}
	return results
}
