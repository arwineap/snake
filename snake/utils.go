package snake

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func checkPointCollision(p1 Point, p2 Point) bool {
	xDistance := math.Abs(float64(p1.X - p2.X))
	yDistance := math.Abs(float64(p1.Y - p2.Y))

	combinedWidth := (float64(p1.Width / 2)) + (float64(p2.Width) / 2)

	if xDistance < combinedWidth && yDistance < combinedWidth {
		return true
	}
	return false
}

func drawPoint(screen *ebiten.Image, p Point, c color.Color) {
	width := float64(p.Width)
	halfWidth := float64(p.Width) / 2
	// It's possible we don't want quarter width and two was just a magic number that matched quarter width
	quarterWidth := float64(p.Width) / 4

	x := float64(p.X)
	y := float64(p.Y)

	ebitenutil.DrawRect(screen, x-halfWidth, y-halfWidth, width, width, color.Black)
	ebitenutil.DrawRect(screen, x-quarterWidth, y-quarterWidth, halfWidth, halfWidth, c)
}

type line struct {
	X1, Y1, X2, Y2 float64
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

type Point struct {
	X     int
	Y     int
	Width int
}

type direction string

const (
	right direction = "right"
	left  direction = "left"
	up    direction = "up"
	down  direction = "down"
)
