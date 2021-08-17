package snake

import "time"

type Option func(s *Snake)

func DefaultOption() []Option {
	return []Option{
		SpeedOption(time.Millisecond * 100),
		DirectionOption(Right),
		WidthOption(16),
	}
}

func SpeedOption(duration time.Duration) Option {
	return func(s *Snake) {
		s.Speed = duration
	}
}

func DirectionOption(direction Direction) Option {
	return func(s *Snake) {
		s.Direction = direction
	}
}

func WidthOption(width int) Option {
	return func(s *Snake) {
		s.Width = width
	}
}

func BorderPaddingOption(padding int) Option {
	return func(s *Snake) {
		s.borderPadding = padding
	}
}
