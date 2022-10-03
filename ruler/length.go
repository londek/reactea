package ruler

// Some LengthNumber - it might be Height, Width or anything else
type (
	LengthNumber    int
	LengthAttribute int
)

type Length interface {
	LengthNumber | LengthAttribute
}

const (
	InitialLength = Auto

	NoContent LengthAttribute = -iota & 0xffffffff
	MinContent
	FitContent
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
