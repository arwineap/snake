package apple

import "time"

type Option func(a *Apple)

func DefaultOption() []Option {
	return []Option{
		FrequencyOption(time.Second),
		WidthOption(16),
	}
}

func FrequencyOption(freq time.Duration) Option {
	return func(a *Apple) {
		a.DropFrequency = freq
	}
}

func WidthOption(width int) Option {
	return func(a *Apple) {
		a.Width = width
	}
}
