package border

import (
	"image/color"
	"snake/point"

	"github.com/hajimehoshi/ebiten/v2"
)

type Border struct {
	Padding int
	Width   int

	positions []point.Point
}

func New(options ...Option) Border {
	if len(options) == 0 {
		options = DefaultOption()
	}

	var b Border
	for _, opt := range options {
		opt(&b)
	}

	return b
}

func (b *Border) Draw(screen *ebiten.Image) {
	if len(b.positions) == 0 {
		w, h := screen.Size()
		b.genBorder(w, h)
	}

	for _, p := range b.positions {
		point.DrawPoint(screen, p, color.Black)
	}
}

func (b *Border) EndGame(head point.Point) (bool, string) {
	for _, borderPoint := range b.positions {
		if point.CheckPointCollision(head, borderPoint) {
			return true, "border collision"
		}
	}
	return false, ""
}

func (b *Border) genBorder(screenWidth int, screenHeight int) {
	var results = []point.Point{{b.Padding, b.Padding, b.Width}}

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
