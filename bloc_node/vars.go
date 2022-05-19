package bloc_node

import (
	"bloc-examples/go/stock/pkg/compare_operator"
	"bloc-examples/go/stock/pkg/raise_or_fall"

	bloc_client "github.com/fBloc/bloc-client-go"
)

var (
	CompareSelections   []bloc_client.SelectOption
	RaiseAndFallOptions []bloc_client.SelectOption
)

func init() {
	compareStrings := compare_operator.AllCompareOperatorStrings()
	for i, j := range compareStrings {
		CompareSelections = append(CompareSelections,
			bloc_client.SelectOption{
				Label: j,
				Value: i,
			},
		)
	}
}

func init() {
	raiseAndFallStrings := raise_or_fall.AllRaiseFallStrings()
	for i, j := range raiseAndFallStrings {
		RaiseAndFallOptions = append(RaiseAndFallOptions,
			bloc_client.SelectOption{
				Label: j,
				Value: i,
			},
		)
	}
}
