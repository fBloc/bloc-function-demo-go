package stock_price_monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"stock/bloc_node"
	"stock/pkg/compare_operator"
	"stock/pkg/realtime_price"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	var _ bloc_client.BlocFunctionNodeInterface = &StockPriceMonitor{}
}

type StockPriceMonitor struct {
}

func (sPM *StockPriceMonitor) AllProgressMilestones() []string {
	return AllMileStones()
}

func (sPM *StockPriceMonitor) IptConfig() bloc_client.Ipts {
	return bloc_client.Ipts{
		{
			Key:     "absolute_price_monitor",
			Display: "absolute_price_monitor",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "stock_code",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
				{
					Hint:            "compare_operator",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   bloc_node.CompareSelections,
					AllowMulti:      false,
				},
				{
					Hint:            "absolute_price",
					ValueType:       bloc_client.FloatValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
			},
		},
	}
}

func (sPM *StockPriceMonitor) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{
		{
			Key:         "suc_msg",
			Description: "suc message",
			ValueType:   bloc_client.StringValueType,
			IsArray:     false,
		},
		{
			Key:         "match_rise",
			Description: "whether input stock match rise monitor",
			ValueType:   bloc_client.BoolValueType,
			IsArray:     false,
		},
		{
			Key:         "match_fall",
			Description: "whether input stock match fall monitor",
			ValueType:   bloc_client.BoolValueType,
			IsArray:     false,
		},
		{
			Key:         "stockCode_map_price",
			Description: "stock code map current price",
			ValueType:   bloc_client.JsonValueType,
			IsArray:     false,
		},
	}
}

func (sPM *StockPriceMonitor) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	logger.Infof("start")

	// need watch stock
	toWatchStockCode, err := ipts.GetStringValue(0, 0)
	if err != nil {
		msg := fmt.Sprintf("get stock code failed: %v", err)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}
	// no stock need to monitor. return
	if toWatchStockCode == "" {
		msg := "no stock needed monitor"
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: true,
			Description:               msg,
		}
		return
	}

	// get compare_operator param
	compareOperatorValue, err := ipts.GetIntValue(0, 1)
	if err != nil {
		msg := fmt.Sprintf("parse compare_operator in needed absolute price monitor failed: %v", err)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}
	compareOperator, err := compare_operator.GetOperatorFromInt(compareOperatorValue)
	if err != nil {
		msg := fmt.Sprintf(
			"compare_operator: %d in needed absolute price monitor not valid",
			compareOperatorValue)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}

	// get absolute_price param
	absolutePrice, err := ipts.GetFloat64Value(0, 2)
	if err != nil {
		msg := fmt.Sprintf(
			"parse absolute_price in needed absolute price monitor failed: %v",
			err)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}
	if absolutePrice < 0 {
		msg := fmt.Sprintf(
			"absolute_price in param should >= 0, but get: %.2f",
			absolutePrice)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}

	// report progress milestone
	progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
		ProgressMilestoneIndex: SucParsedParam.MilestoneIndex(),
	}
	progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
		ProgressMilestoneIndex: StartVisitRemoteApi4StockRealtimePrice.MilestoneIndex(),
	}

	var sucMsg string
	var matchRise, matchFall bool
	stockCodeMapPrice := make(map[string]float64)
	// param valid, going to check price
	price, err := realtime_price.GetRealtimePrice(toWatchStockCode)
	if err != nil {
		msg := fmt.Sprintf(
			"stock_code: %s visit remote to get realtime price failed: %v",
			toWatchStockCode, err)
		logger.Errorf(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  msg,
		}
		return
	}
	stockCodeMapPrice[toWatchStockCode] = price
	progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
		ProgressMilestoneIndex: FinishedVisitRemoteApi4StockRealtimePrice.MilestoneIndex(),
	}

	stockCodeMapPriceByte, _ := json.Marshal(stockCodeMapPrice)
	// hit true
	if compare_operator.CompareFloat64(price, compareOperator, absolutePrice) {
		logger.Infof("hit absolute price monitor")
		sucMsg += fmt.Sprintf(
			"%s is %s %.2f",
			toWatchStockCode, compareOperator.String(), absolutePrice)

		if compareOperator == compare_operator.Gt ||
			compareOperator == compare_operator.Gte {
			matchRise = true
		} else if compareOperator == compare_operator.Lt ||
			compareOperator == compare_operator.Lte {
			matchFall = true
		}
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: false,
			Description:               sucMsg,
			// Detail should mathc the OptConfig
			Detail: map[string]interface{}{
				"suc_msg":             sucMsg,
				"match_rise":          matchRise,
				"match_fall":          matchFall,
				"stockCode_map_price": stockCodeMapPriceByte, // json type should return data after serialize
			},
		}
		return
	}
	logger.Infof("miss absolute price monitor")
	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc:                       true,
		InterceptBelowFunctionRun: true,
		Description:               sucMsg,
		// Detail should mathc the OptConfig
		Detail: map[string]interface{}{
			"suc_msg":             sucMsg,
			"match_rise":          matchRise,
			"match_fall":          matchFall,
			"stockCode_map_price": stockCodeMapPriceByte, // json type should return data after serialize
		},
	}
}
