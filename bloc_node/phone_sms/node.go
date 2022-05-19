package phone_sms

import (
	"bloc-examples/go/stock/pkg/sms"
	"context"
	"encoding/json"
	"fmt"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	// should satisfy the bloc_node's interface
	var _ bloc_client.BlocFunctionNodeInterface = &SMS{}
}

type SMS struct {
}

// AllProgressMilestones define the progress milestones of this function,
// as this function is not a long run function, choose to not set it.
func (s *SMS) AllProgressMilestones() []string {
	return []string{}
}

// IptConfig define the ipt param config of this function
func (s *SMS) IptConfig() bloc_client.Ipts {
	return bloc_client.Ipts{
		{
			Key:     "content",
			Display: "content",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "msg",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
			},
		},
		{
			Key:     "phone_numbers",
			Display: "phone_numbers",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "phone_numbers",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      true,
				},
			},
		},
	}
}

// OptConfig define the opt config of this function
func (s *SMS) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{
		{
			Key:         "suc_numbers",
			Description: "suc send phone numbers",
			ValueType:   bloc_client.StringValueType,
			IsArray:     true,
		},
		{
			Key:         "fail_number_map_fail_msg",
			Description: "fail sended phone number and its fail reason msg",
			ValueType:   bloc_client.JsonValueType,
			IsArray:     false,
		},
	}
}

// Run function's actual execute logic
func (s *SMS) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	logger.Infof("start")

	// get msg content param
	content, err := ipts.GetStringValue(0, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get msg content from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
		}
		return
	}
	if content == "" {
		msg := "msg content is empty, no need to send"
		logger.Infof(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: true,
			Description:               msg,
		}
		return
	}
	logger.Infof("get msg content suc")

	// get to send phone numbers param
	phoneNumbers, err := ipts.GetStringSliceValue(1, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get phone number list from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
		}
		return
	}
	if len(phoneNumbers) == 0 {
		msg := "phone list is empty, no need to send"
		logger.Infof(msg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: true,
			Description:               msg,
		}
		return
	}
	logger.Infof("get %d phone numbers", len(phoneNumbers))

	// do the send logic
	sucNumbers := make([]string, 0, len(phoneNumbers))
	failNumberMapFailMsg := make(map[string]string)

	for _, phoneNumber := range phoneNumbers {
		err := sms.SendMsg(phoneNumber, content)
		if err == nil {
			logger.Infof("phone %s send suc", phoneNumber)
			sucNumbers = append(sucNumbers, phoneNumber)
			continue
		}
		logger.Errorf("phone %s send failed, msg: %v", phoneNumber, err)
		failNumberMapFailMsg[phoneNumber] = err.Error()
	}
	logger.Infof("finish send")

	failNumberMapFailMsgByte, err := json.Marshal(failNumberMapFailMsg)
	if err != nil {
		logger.Errorf("marshal failNumberMapFailMsg failed: %v", err)
	}

	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc:                       true,
		InterceptBelowFunctionRun: false,
		Description: fmt.Sprintf(
			"suc send %d phone numbers, fail send %d",
			len(sucNumbers), len(failNumberMapFailMsg)),
		Detail: map[string]interface{}{
			"suc_numbers":              sucNumbers,
			"fail_number_map_fail_msg": failNumberMapFailMsgByte,
		},
	}
}
