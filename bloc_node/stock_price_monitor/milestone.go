package stock_price_monitor

type progressMileStone int

const (
	SucParsedParam progressMileStone = iota
	StartVisitRemoteApi4StockRealtimePrice
	FinishedVisitRemoteApi4StockRealtimePrice
	maxMilestone
)

func (pMS progressMileStone) String() string {
	switch pMS {
	case SucParsedParam:
		return "parse param suc"
	case StartVisitRemoteApi4StockRealtimePrice:
		return "start visit remote api 4 stock realtime price"
	case FinishedVisitRemoteApi4StockRealtimePrice:
		return "finished visit remote api 4 stock realtime price"
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
