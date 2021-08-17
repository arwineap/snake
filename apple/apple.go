package apple

import (
	"image/color"
	"snake/point"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Apple struct {
	Positions     []point.Point
	LastDrop      time.Time
	DropFrequency time.Duration
	Width         int
	extendQueue   int
}

func New(options ...Option) *Apple {
	var result Apple

	if len(options) == 0 {
		options = DefaultOption()
	}

	for _, opt := range options {
		opt(&result)
	}
	return &result
}

func (a *Apple) Reset() {
	a.Positions = []point.Point{}
}

func (a *Apple) ExtendSnake(head point.Point) int {
	collisions := a.Collisions(head)
	lenCollisions := len(collisions)

	for _, c := range collisions {
		a.Remove(c)
	}

	return lenCollisions

}

func (a *Apple) Drop(position point.Point) {
	if time.Since(a.LastDrop) < a.DropFrequency {
		return
	}
	position.Width = a.Width
	a.Positions = append(a.Positions, position)
	a.LastDrop = time.Now()
}

func (a *Apple) Remove(p point.Point) {
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
		point.DrawPoint(screen, p, color.RGBA{R: 255, G: 100, B: 100, A: 255})
	}
}

// PointCollides - check if any apple collides with point
func (a *Apple) PointCollides(point point.Point) bool {
	appleCollisions := a.Collisions(point)
	if len(appleCollisions) == 0 {
		return false
	}
	return true
}

// Collisions - returns all points that collide
func (a *Apple) Collisions(pnt point.Point) []point.Point {
	var results []point.Point
	for _, p := range a.Positions {
		if point.CheckPointCollision(p, pnt) {
			results = append(results, p)
		}
	}
	return results
}
