package snake

import (
	"image/color"
	"snake/point"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Snake struct {
	Position      []point.Point
	Direction     Direction
	Speed         time.Duration
	LastMove      time.Time
	Width         int
	extendQueue   int
	borderPadding int
}

func New(options ...Option) *Snake {
	if len(options) == 0 {
		options = DefaultOption()
	}

	var result Snake
	for _, option := range options {
		option(&result)
	}

	return &result
}

func (s *Snake) Reset() {
	s.Position = []point.Point{{(s.Width * 3) + s.borderPadding, (s.borderPadding + 1) + (s.Width / 2), s.Width}, {(s.Width * 2) + s.borderPadding, (s.borderPadding + 1) + (s.Width / 2), s.Width}}
	s.Direction = Right
}

func (s *Snake) Draw(screen *ebiten.Image) {
	for _, p := range s.Position {
		point.DrawPoint(screen, p, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	}
}

func (s *Snake) NextPoint() point.Point {
	var newPoint point.Point
	currentPoint := s.Head()
	switch s.Direction {
	case Up:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y - s.Width
	case Down:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y + s.Width
	case Left:
		newPoint = currentPoint
		newPoint.X = newPoint.X - s.Width
	case Right:
		newPoint = currentPoint
		newPoint.X = newPoint.X + s.Width
	}

	return newPoint
}

func (s *Snake) EndGame(head point.Point) (bool, string) {
	// This needs to check if we collide with any point in the snake except the head

	for i, p := range s.Position {
		if i != 0 {
			if point.CheckPointCollision(head, p) {
				return true, "snake collision"
			}
		}
	}
	return false, ""
}

// PointCollides - check if any unit of snake collides with point
func (s *Snake) PointCollides(pnt point.Point) bool {
	for _, p := range s.Position {
		if point.CheckPointCollision(p, pnt) {
			return true
		}
	}
	return false
}

// Head - return the point address of the snake's head
func (s *Snake) Head() point.Point {
	return s.Position[0]
}

func (s *Snake) Move() {
	if time.Since(s.LastMove) < s.Speed {
		return
	}

	defer func() { s.LastMove = time.Now() }()
	// To move the snake we prepend a Point to the start of the slice, then remove the final Position of the slice
	// If we moved into an apple, we don't remove the final Position
	var newPoint point.Point
	head := s.Head()
	switch s.Direction {
	case Up:
		newPoint = point.Point{X: head.X, Y: head.Y - s.Width, Width: s.Width}
	case Down:
		newPoint = point.Point{X: head.X, Y: head.Y + s.Width, Width: s.Width}
	case Left:
		newPoint = point.Point{X: head.X - s.Width, Y: head.Y, Width: s.Width}
	case Right:
		newPoint = point.Point{X: head.X + s.Width, Y: head.Y, Width: s.Width}
	}

	s.Position = append([]point.Point{newPoint}, s.Position...)

	if s.extendQueue == 0 {
		// don't extendQueue, which means we need to remove the last position so snake "moves"
		s.Position = s.Position[:len(s.Position)-1]
		return
	}
	s.extendQueue--
}

func (s *Snake) AddExtend(x int) {
	s.extendQueue = s.extendQueue + x
}

func (s *Snake) SetDirection(d Direction) {
	s.Direction = d
}
