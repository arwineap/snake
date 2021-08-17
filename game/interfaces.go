package game

import (
	"snake/point"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawer interface {
	Draw(*ebiten.Image)
}

type Extender interface {
	ExtendSnake(head point.Point) int
}

type Ender interface {
	EndGame(head point.Point) (bool, string)
}

type Reseter interface {
	Reset()
}

type Droper interface {
	Drop(unusedPoint point.Point)
}
