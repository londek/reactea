package render

import "fmt"

type (
	MarginNumber     = int
	MarginPercentage = float32
	MarginAttribute  int
)

type marginType int

var (
	zeroMarginValue = marginValue{marginNumberType, 0}
	zeroMargin      = margin{zeroMarginValue, zeroMarginValue, zeroMarginValue, zeroMarginValue}
)

const (
	marginNumberType marginType = iota
	marginPercentageType
	marginAttributeType
)

func (mt marginType) String() string {
	switch mt {
	case marginNumberType:
		return "MarginNumber"
	case marginPercentageType:
		return "MarginPercentage"
	case marginAttributeType:
		return "MarginAttribute"
	default:
		return "<unknown>"
	}
}

type margin struct {
	Top, Right, Bottom, Left marginValue
}

func (m margin) String() string {
	return fmt.Sprintf("(top: %s - right: %s - bottom: %s - left: %s)", m.Top, m.Right, m.Bottom, m.Left)
}

func Margin(values ...any) margin {
	if len(values) == 1 {
		all := MarginValue(values[0])
		return margin{all, all, all, all}
	} else if len(values) == 2 {
		vertical, horizontal := MarginValue(values[0]), MarginValue(values[1])
		return margin{vertical, horizontal, vertical, horizontal}
	} else if len(values) == 3 {
		top, horizontal, bottom := MarginValue(values[0]), MarginValue(values[1]), MarginValue(values[2])
		return margin{top, horizontal, bottom, horizontal}
	} else if len(values) >= 4 {
		top, right, bottom, left := MarginValue(values[0]), MarginValue(values[1]), MarginValue(values[2]), MarginValue(values[3])
		return margin{top, right, bottom, left}
	} else {
		return zeroMargin
	}
}

func MarginValue(value any) marginValue {
	switch value := any(value).(type) {
	case MarginNumber:
		return marginValue{marginNumberType, value}
	case MarginPercentage:
		return marginValue{marginPercentageType, int(value * floatMultiplier)}
	case MarginAttribute:
		return marginValue{marginAttributeType, int(value)}
	default:
		panic("invalid marginValue type (should never happen)")
	}
}

type marginValue struct {
	marginType marginType
	underlying int
}

func (mv marginValue) String() string {
	switch mv.marginType {
	case marginNumberType:
		return fmt.Sprintf("%d", mv.underlying)
	case marginPercentageType:
		return fmt.Sprintf("%f%%", float32(mv.underlying)/floatMultiplier)
	case marginAttributeType:
		return LengthAttribute(mv.underlying).String()
	default:
		return "<unknown>"
	}
}

const (
	AutoMargin MarginAttribute = -iota & 0xffffffff
)

func (ma MarginAttribute) String() string {
	switch ma {
	case AutoMargin:
		return "Auto"
	default:
		return "<unknown>"
	}
}
