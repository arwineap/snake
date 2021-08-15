package snake

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Snake struct {
	Position    []Point
	Direction   direction
	Speed       time.Duration
	LastMove    time.Time
	Width       int
	extendQueue int
}

func (s *Snake) Draw(screen *ebiten.Image) {
	for _, p := range s.Position {
		drawPoint(screen, p, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	}
}

func (s *Snake) NextPoint() Point {
	var newPoint Point
	currentPoint := s.Head()
	switch s.Direction {
	case up:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y - s.Width
	case down:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y + s.Width
	case left:
		newPoint = currentPoint
		newPoint.X = newPoint.X - s.Width
	case right:
		newPoint = currentPoint
		newPoint.X = newPoint.X + s.Width
	}

	return newPoint
}

// PointCollides - check if any unit of snake collides with point
func (s *Snake) PointCollides(point Point) bool {
	for _, p := range s.Position {
		if checkPointCollision(p, point) {
			return true
		}
	}
	return false
}

// Head - return the point address of the snake's head
func (s *Snake) Head() Point {
	return s.Position[0]
}

func (s *Snake) Move() bool {
	if time.Since(s.LastMove) < s.Speed {
		return false
	}

	// To move the Snake we prepend a Point to the start of the slice, then remove the final Position of the slice
	// If we moved into an apple, we don't remove the final Position
	var newPoint Point
	head := s.Head()
	switch s.Direction {
	case up:
		newPoint = Point{X: head.X, Y: head.Y - s.Width, Width: s.Width}
	case down:
		newPoint = Point{X: head.X, Y: head.Y + s.Width, Width: s.Width}
	case left:
		newPoint = Point{X: head.X - s.Width, Y: head.Y, Width: s.Width}
	case right:
		newPoint = Point{X: head.X + s.Width, Y: head.Y, Width: s.Width}
	default:
		return true
	}

	s.Position = append([]Point{newPoint}, s.Position...)

	s.LastMove = time.Now()
	return true
}

func (s *Snake) AddExtend(x int) {
	s.extendQueue = s.extendQueue + x
}

func (s *Snake) Extend() {
	if s.extendQueue == 0 {
		// don't extendQueue, which means we need to remove the last position so snake "moves"
		s.Position = s.Position[:len(s.Position)-1]
		return
	}
	s.extendQueue--
}

func (s *Snake) SetDirection(d direction) {
	s.Direction = d
}
