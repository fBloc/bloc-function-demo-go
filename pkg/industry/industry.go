package industry

type Industry int

const (
	Agriculture = iota
	Alcohol
	IT
	PublicTransportation
	Semimonomer
	Chemical
	Restaurant
	Aviation
	AutoParts
	max // just a few listed
)

func (i Industry) String() string {
	switch i {
	case Agriculture:
		return "agriculture"
	case IT:
		return "it"
	case PublicTransportation:
		return "public_transportation"
	case Semimonomer:
		return "semimonomer"
	case Chemical:
		return "chemical"
	case Restaurant:
		return "restaurant"
	case Aviation:
		return "aviation"
	case AutoParts:
		return "auto_parts"
	default:
		return "unknown"
	}
}

func AllIndustryStrings() []string {
	ret := make([]string, 0, max-1)
	for i := Agriculture; i < max; i++ {
		ret = append(ret, Industry(i).String())
	}
	return ret
}
