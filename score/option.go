package score

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Option func(*Score) error

func DefaultOption() []Option {
	return []Option{
		CountOption(0),
		FontOption(16, 72),
	}
}

func CountOption(count int) Option {
	return func(s *Score) error {
		s.Count = count
		return nil
	}
}

func FontOption(size float64, dpi float64) Option {
	tt, ttErr := opentype.Parse(fonts.MPlus1pRegular_ttf)
	fnt, fntErr := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	return func(s *Score) error {
		if ttErr != nil {
			return ttErr
		}
		if fntErr != nil {
			return fntErr
		}
		s.Font = fnt
		return nil
	}
}
