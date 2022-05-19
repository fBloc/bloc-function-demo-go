package sleep

type progressMileStone int

const (
	SucParsedParam progressMileStone = iota
	StartSleep
	FinishedSleep
	maxMilestone
)

func (pMS progressMileStone) String() string {
	switch pMS {
	case SucParsedParam:
		return "parse param suc"
	case StartSleep:
		return "start sleep"
	case FinishedSleep:
		return "finished sleep"
	}
	return "unknown"
}

func (pMS progressMileStone) MilestoneIndex() *int {
	tmp := int(pMS)
	return &tmp
}

func AllMileStones() []string {
	ret := make([]string, 0, maxMilestone-1)
	for i := SucParsedParam; i < maxMilestone; i++ {
		ret = append(ret, progressMileStone(i).String())
	}
	return ret
}
