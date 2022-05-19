package compare_operator

import "fmt"

type CompareOperator int

const (
	Eq CompareOperator = iota
	Gt
	Gte
	Lt
	Lte
	max
)

func (cO CompareOperator) String() string {
	switch cO {
	case Eq:
		return "="
	case Gt:
		return ">"
	case Gte:
		return "≥"
	case Lt:
		return "<"
	case Lte:
		return "≤"
	default:
		return "unknown"
	}
}

func AllCompareOperatorStrings() []string {
	ret := make([]string, 0, max-1)
	for i := Eq; i < max; i++ {
		ret = append(ret, CompareOperator(i).String())
	}
	return ret
}

func GetOperatorFromInt(i int) (CompareOperator, error) {
	if i >= int(max) {
		return max, fmt.Errorf("invalid value: %d", i)
	}
	return CompareOperator(i), nil
}

func CompareFloat64(a float64, cO CompareOperator, b float64) bool {
	switch cO {
	case Eq:
		return a == b
	case Gt:
		return a > b
	case Gte:
		return a >= b
	case Lt:
		return a < b
	case Lte:
		return a <= b
	default:
		return false
	}
}
