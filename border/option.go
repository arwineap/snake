package border

type Option func(*Border)

func DefaultOption() []Option {
	return []Option{
		PaddingOption(20),
		WidthOption(2),
	}
}

func WidthOption(width int) Option {
	return func(b *Border) {
		b.Width = width
	}
}

func PaddingOption(padding int) Option {
	return func(b *Border) {
		b.Padding = padding
	}
}
