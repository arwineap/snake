package snake

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Border struct {
	Padding int
	Width   int

	positions []Point
}

func (b *Border) Draw(screen *ebiten.Image) {
	if len(b.positions) == 0 {
		w, h := screen.Size()
		b.genBorder(w, h)
	}

	for _, p := range b.positions {
		drawPoint(screen, p, color.Black)
	}
}

func (b *Border) Collides(p Point) bool {
	for _, borderPoint := range b.positions {
		if checkPointCollision(p, borderPoint) {
			return true
		}
	}
	return false
}

func (b *Border) genBorder(screenWidth int, screenHeight int) {
	var results = []Point{{b.Padding, b.Padding, b.Width}}

	var newPoint = results[0]
	// Top points
	for newPoint.X < (screenWidth - b.Padding) {
		newPoint.X = newPoint.X + b.Width
		results = append(results, newPoint)
	}

	// Right wall points
	for newPoint.Y <= (screenHeight - b.Padding) {
		newPoint.Y = newPoint.Y + b.Width
		results = append(results, newPoint)
	}

	// Bottom wall points
	for newPoint.X >= b.Padding {
		newPoint.X = newPoint.X - b.Width
		results = append(results, newPoint)
	}

	// Left wall points
	for newPoint.Y > b.Padding {
		newPoint.Y = newPoint.Y - b.Width
		results = append(results, newPoint)
	}

	b.positions = results
}
