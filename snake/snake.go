package snake

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Snake struct {
	Position  []Point
	Direction direction
	Speed     time.Duration
	LastMove  time.Time
	Width     int
}

func (s Snake) Draw(screen *ebiten.Image) {
	for _, p := range s.Position {
		drawPoint(screen, p, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	}
}
