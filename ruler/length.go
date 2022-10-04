package ruler

// Some LengthNumber - it might be Height, Width or anything else
type (
	LengthNumber     int
	LengthPercentage float32
	LengthAttribute  int
)

type Length interface {
	LengthNumber | LengthPercentage | LengthAttribute
}

func (l LengthNumber) Get(*RenderContext) int {
	return int(l)
}

const (
	// Ignore attribute, use numbers
	NoAttribute LengthAttribute = -iota & 0xffffffff
	MinContent
	FitContent // MaxContent <= FitContent <= MinContent
	MaxContent
	Auto
)

func (l *LengthNumber) Take(n LengthNumber) LengthNumber {
	if *l >= n {
		*l -= n
		return n
	}

	taken := *l
	*l = 0

	return taken
}

func (l LengthNumber) Left() LengthNumber {
	return l
}
