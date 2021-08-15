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

	if xDistance < 8 && yDistance < 8 {
		return true
	}
	return false
}

func drawPoint(screen *ebiten.Image, p Point, c color.Color) {
	ebitenutil.DrawRect(screen, float64(p.X)-4, float64(p.Y)-4, float64(8), float64(8), color.Black)
	ebitenutil.DrawRect(screen, float64(p.X)-2, float64(p.Y)-2, float64(4), float64(4), c)
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
	X int
	Y int
}

type direction string

const (
	right direction = "right"
	left  direction = "left"
	up    direction = "up"
	down  direction = "down"
)
