package new_stock

import (
	"encoding/json"
	"testing"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func TestNewStockWithNoFilter(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&NewStock{},
		[][]interface{}{
			{
				[]string{},
			},
			{
				[]int{},
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	returnedStockCodes := executeOpt.Detail["new_stock_codes"]
	returnedStockCodesSlice, ok := returnedStockCodes.([]string)
	if !ok {
		t.Fatalf("opt detail new_stock_codes field should be []string type, convert failed")
	}
	if len(returnedStockCodesSlice) != len(fakeNewStock) {
		t.Fatal("amount not match")
	}
}

func TestNewStockWithFilter(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&NewStock{},
		[][]interface{}{
			{
				[]string{fakeNewStock[0].Exchange},
			},
			{
				[]int{},
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	// opt new_stock_codes slice field
	returnedStockCodes := executeOpt.Detail["new_stock_codes"]
	returnedStockCodesSlice, ok := returnedStockCodes.([]string)
	if !ok {
		t.Fatalf("opt detail new_stock_codes field should be []string type, convert failed")
	}
	if len(returnedStockCodesSlice) != 1 {
		t.Fatal("new_stock_codes amount not match")
	}
	// opt new_stocks json field
	newStocksField := executeOpt.Detail["new_stocks"]
	newStocksByte, ok := newStocksField.([]byte)
	if !ok {
		t.Fatal("new_stocks to  failed")
	}

	var newStocks NewStocks
	err := json.Unmarshal(newStocksByte, &newStocks)
	if err != nil {
		t.Fatalf("json unmarshal new_stocks failed with %s", err.Error())
	}
	if len(newStocks) != 1 {
		t.Fatal("new_stocks amount not match")
	}
	if newStocks[0].Exchange != fakeNewStock[0].Exchange {
		t.Fatal("new_stocks exchange not match")
	}
}
