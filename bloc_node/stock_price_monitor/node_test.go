package stock_price_monitor

import (
	"testing"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func TestStockPriceMonitorSuc(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&StockPriceMonitor{},
		[][]interface{}{
			{
				"001235.SZ", // stock_code
				2,           // compare_operator - ≥
				0.0,         // absolute_price - 0 to make sure it's 100% hit
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should suc, but fail.")
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept_below_functionrun")
	}
	matchRisedStockCodesInterface := executeOpt.Detail["match_rise_stock_codes"]
	matchRisedStockCodes, ok := matchRisedStockCodesInterface.([]string)
	if !ok {
		t.Fatalf("opt detail match_rise_stock_codes field should be []string type, convert failed")
	}
	if len(matchRisedStockCodes) != 1 {
		t.Fatal("amount not match")
	}
}
