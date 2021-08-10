package snake

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Score struct {
	Count int
	Font  font.Face
}

func (s Score) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("score: %d", s.Count), s.Font, 180, 10, color.White)
}
