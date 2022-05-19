package stock_exchange

type StockExchange int

const (
	ShanghaiExchange StockExchange = iota
	ShenzhenExchange
	maxStockExchange
)

const notValidStockExchangeString = "unknown"

func (sE StockExchange) String() string {
	switch sE {
	case ShanghaiExchange:
		return "shanghai"
	case ShenzhenExchange:
		return "shenzhen"
	default:
		return "unknown"
	}
}

func (sE StockExchange) ExchangeCode() string {
	switch sE {
	case ShanghaiExchange:
		return "SSE"
	case ShenzhenExchange:
		return "SZSE"
	default:
		return "unknown"
	}
}

func AllExchangeNameAndCode() [][]string {
	ret := make([][]string, 0, maxStockExchange-1)
	for i := ShanghaiExchange; i < maxStockExchange; i++ {
		ret = append(ret, []string{i.String(), i.ExchangeCode()})
	}
	return ret
}
