package raise_or_fall

type RaiseAndFall int

const (
	Raise RaiseAndFall = iota
	Fall
	max
)

func (rAF RaiseAndFall) String() string {
	switch rAF {
	case Raise:
		return "raise"
	case Fall:
		return "fall"
	default:
		return "unknown"
	}
}

func AllRaiseFallStrings() []string {
	ret := make([]string, 0, max-1)
	for i := Raise; i < max; i++ {
		ret = append(ret, RaiseAndFall(i).String())
	}
	return ret
}
