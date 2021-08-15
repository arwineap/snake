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
	w, h := screen.Size()

	text.Draw(screen, fmt.Sprintf("score: %d", s.Count), s.Font, int(float64(w)*0.90), int(float64(h)*0.015), color.White)
}
