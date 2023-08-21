package render

import "fmt"

type (
	LengthNumber     = int
	LengthPercentage = float32
	LengthAttribute  int
)

var (
	defaultLength = length{lengthAttributeType, int(Auto)}
	zeroLength    = length{lengthNumberType, 0}
)

type lengthType int

const (
	lengthNumberType lengthType = iota
	lengthPercentageType
	lengthAttributeType
)

func (lt lengthType) String() string {
	switch lt {
	case lengthNumberType:
		return "LengthNumber"
	case lengthPercentageType:
		return "LengthPercentage"
	case lengthAttributeType:
		return "LengthAttribute"
	default:
		return "<unknown>"
	}
}

type length struct {
	lengthType lengthType
	underlying int
}

func (l length) Is(value any) bool {
	if value, ok := value.(length); ok {
		return l.lengthType == value.lengthType && l.underlying == value.underlying
	}

	return l.underlying == value
}

func (l length) String() string {
	switch l.lengthType {
	case lengthNumberType:
		return fmt.Sprintf("%d", l.underlying)
	case lengthPercentageType:
		return fmt.Sprintf("%f%%", float32(l.underlying)/floatMultiplier)
	case lengthAttributeType:
		return LengthAttribute(l.underlying).String()
	default:
		return "<unknown>"
	}
}

func Length(value any) length {
	switch value := any(value).(type) {
	case LengthNumber:
		return length{lengthNumberType, value}
	case LengthPercentage:
		return length{lengthPercentageType, int(value * floatMultiplier)}
	case LengthAttribute:
		return length{lengthAttributeType, int(value)}
	default:
		panic("invalid length type (should never happen)")
	}
}

const (
	Auto LengthAttribute = -iota & 0xffffffff
	MinContent
	FitContent // MaxContent <= FitContent <= MinContent
	MaxContent
)

func (la LengthAttribute) String() string {
	switch la {
	case Auto:
		return "Auto"
	case MinContent:
		return "MinContent"
	case FitContent:
		return "FitContent"
	case MaxContent:
		return "MaxContent"
	default:
		return "<unknown>"
	}
}

// func (l *LengthNumber) Take(n LengthNumber) LengthNumber {
// 	if *l >= n {
// 		*l -= n
// 		return n
// 	}

// 	taken := *l
// 	*l = 0

// 	return taken
// }

// func (l LengthNumber) Left() LengthNumber {
// 	return l
// }
