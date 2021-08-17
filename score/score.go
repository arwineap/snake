package score

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

func New(options ...Option) (*Score, error) {
	if len(options) == 0 {
		options = DefaultOption()
	}

	var s Score
	for _, opt := range options {
		err := opt(&s)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s Score) Draw(screen *ebiten.Image) {
	w, h := screen.Size()

	text.Draw(screen, fmt.Sprintf("score: %d", s.Count), s.Font, int(float64(w)*0.90), int(float64(h)*0.015), color.White)
}

func (s *Score) Reset() {
	s.Count = 0
}