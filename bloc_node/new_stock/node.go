package new_stock

import (
	"bloc-examples/go/stock/pkg/industry"
	"bloc-examples/go/stock/pkg/stock_exchange"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	var _ bloc_client.BlocFunctionNodeInterface = &NewStock{}
}

type NewStock struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Price       float32 `json:"price"`
	IssueAmount uint64  `json:"issue_amount"`
	IssueDate   string  `json:"issue_date"`
	Exchange    string  `json:"exchange"`
	Industry    string  `json:"industry"`
}

type NewStocks []NewStock

func (nS NewStocks) json() ([]byte, error) {
	jsonBytes, err := json.Marshal(nS)
	if err != nil {
		return []byte{}, err
	}
	return jsonBytes, nil
}

func (nS NewStocks) String() string {
	msgs := make([]string, 0, len(nS))
	for _, i := range nS {
		msgs = append(
			msgs,
			fmt.Sprintf(
				"new stock: name-%s, code-%s, price-%.2f, issue_amount-%d, issue_date-%s, exchange-%s",
				i.Name, i.Code, i.Price, i.IssueAmount, i.IssueDate, i.Exchange,
			),
		)
	}
	return strings.Join(msgs, ";")
}

func (nStock *NewStock) AllProgressMilestones() []string {
	// as this is a run fast function, choose not to report progress milestone
	return []string{}
}

func (nStock *NewStock) IptConfig() bloc_client.Ipts {
	// all exchange options
	exchangeNameAndCode := stock_exchange.AllExchangeNameAndCode()
	exchangeSelectOptions := make([]bloc_client.SelectOption, 0, len(exchangeNameAndCode))
	for _, exchange := range exchangeNameAndCode {
		exchangeSelectOptions = append(exchangeSelectOptions, bloc_client.SelectOption{
			Label: exchange[0],
			Value: exchange[1],
		})
	}

	// all industry options
	industryStrings := industry.AllIndustryStrings()
	industryOptions := make([]bloc_client.SelectOption, 0, len(industryStrings))
	for i, j := range industryStrings {
		industryOptions = append(industryOptions, bloc_client.SelectOption{
			Label: j,
			Value: i,
		})
	}

	return bloc_client.Ipts{
		// filter stock exchange param
		{
			Key:     "exchange",
			Display: "filter certain stock exchange",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "exchange",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   exchangeSelectOptions,
					AllowMulti:      true,
				},
			},
		},
		// filter stock industry param
		{
			Key:     "industry",
			Display: "filter stock industry",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "industry",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   industryOptions,
					AllowMulti:      true,
				},
			},
		},
	}
}

func (nStock *NewStock) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{
		{
			Key:         "error_msg",
			Description: "error message",
			ValueType:   bloc_client.StringValueType,
			IsArray:     false,
		},
		{
			Key:         "suc_msg",
			Description: "success message",
			ValueType:   bloc_client.StringValueType,
			IsArray:     false,
		},
		{
			Key:         "new_stock_codes",
			Description: "match filter new stock's code array",
			ValueType:   bloc_client.StringValueType,
			IsArray:     true,
		},
		{
			Key:         "new_stocks",
			Description: "match filter new stock array",
			ValueType:   bloc_client.JsonValueType,
			IsArray:     false,
		},
	}
}

func (nStock *NewStock) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	logger.Infof("start")

	exchangeSlice, err := ipts.GetStringSliceValue(0, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get exchange slice from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
			// Detail should mathc the OptConfig
			Detail: map[string]interface{}{
				"error_msg":  errorMsg,
				"new_stocks": "",
			},
		}
		return
	}
	logger.Infof("parse exchange_slice from ipt suc")

	exchangeMap := make(map[string]bool, len(exchangeSlice))
	if len(exchangeSlice) > 0 {
		for _, i := range exchangeSlice {
			exchangeMap[i] = true
		}
	}

	industrySlice, err := ipts.GetIntSliceValue(1, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get industry slice from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
			// Detail should mathc the OptConfig
			Detail: map[string]interface{}{
				"error_msg":  errorMsg,
				"new_stocks": "",
			},
		}
		return
	}

	industryMap := make(map[string]bool, len(exchangeSlice))
	if len(exchangeSlice) > 0 {
		for _, i := range industrySlice {
			industryMap[industry.Industry(i).String()] = true
		}
	}
	logger.Infof("parse industry_slice from ipt suc")

	logger.Infof("start filter stocks")
	var hitNewStocks NewStocks
	var hitStockIDs []string
	for _, i := range fakeNewStock {
		exchangeHit := len(exchangeSlice) == 0
		if _, ok := exchangeMap[i.Exchange]; ok {
			exchangeHit = true
		}

		industryHit := len(industrySlice) == 0
		if _, ok := industryMap[i.Industry]; ok {
			industryHit = true
		}

		if exchangeHit && industryHit {
			hitNewStocks = append(hitNewStocks, i)
			hitStockIDs = append(hitStockIDs, i.Code)
		}
	}
	msg := fmt.Sprintf(
		"filter stocks finished with ipt amount: %d, opt amount: %d",
		len(fakeNewStock), len(hitNewStocks))
	logger.Infof(msg)

	if len(hitStockIDs) == 0 {
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: true,
			Description:               msg,
		}
		return
	}

	hitNewStocksJson, err := hitNewStocks.json()
	if err != nil {
		errorMsg := fmt.Sprintf("json encode new_stocks failed: %v", err)
		logger.Errorf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
			// Detail should mathc the OptConfig
			Detail: map[string]interface{}{
				"error_msg":       errorMsg,
				"new_stocks":      "",
				"new_stock_codes": hitStockIDs,
			},
		}
		return
	}

	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc:                       true,
		InterceptBelowFunctionRun: false,
		Description:               msg,
		// Detail should mathc the OptConfig
		Detail: map[string]interface{}{
			"error_msg":       "",
			"suc_msg":         hitNewStocks.String(),
			"new_stock_codes": hitStockIDs,
			"new_stocks":      hitNewStocksJson,
		},
	}
}

func TodayString() string {
	return time.Now().Format("20060102")
}

var fakeNewStock = NewStocks{
	{
		Name:        "里得电科",
		Code:        "001235.SZ",
		Price:       25.48,
		IssueAmount: 2121,
		IssueDate:   TodayString(), // today issue(fake)
		Exchange:    "SSE",
		Industry:    "it",
	},
	{
		Name:        "荣亿精密",
		Code:        "873223.BJ",
		Price:       3.21,
		IssueAmount: 3790,
		IssueDate:   TodayString(), // today issue(fake)
		Exchange:    "SZSE",
		Industry:    "auto_parts",
	},
}
