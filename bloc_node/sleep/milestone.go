package sleep

type progressMileStone int

const (
	ParsingParam progressMileStone = iota
	Sleeping
	Finish
	maxMilestone
)

func (pMS progressMileStone) String() string {
	switch pMS {
	case ParsingParam:
		return "parsing param"
	case Sleeping:
		return "sleeping"
	case Finish:
		return "finish"
	}
	return "unknown"
}

func (pMS progressMileStone) MilestoneIndex() *int {
	tmp := int(pMS)
	return &tmp
}

func AllMileStones() []string {
	ret := make([]string, 0, maxMilestone-1)
	for i := ParsingParam; i < maxMilestone; i++ {
		ret = append(ret, progressMileStone(i).String())
	}
	return ret
}
