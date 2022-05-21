/*
   The sleep package implements a sleep function node to bloc

   this function node is simulate a long run node,
   demonstration how to real-time reporting progress_milestone & progress percent,
   which can be seen & dynamic update in the bloc-frontend UI while running.
*/
package sleep

import (
	"context"
	"fmt"
	"time"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	var _ bloc_client.BlocFunctionNodeInterface = &Sleep{}
}

type Sleep struct {
}

func (s *Sleep) AllProgressMilestones() []string {
	return AllMileStones()
}

func (s *Sleep) IptConfig() bloc_client.Ipts {
	return bloc_client.Ipts{
		{
			Key:     "sleep",
			Display: "sleep",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "sleep_seconds",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
			},
		},
	}
}

func (s *Sleep) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{}
}

func (s *Sleep) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	logger.Infof("start") // report log

	progressReportChan <- bloc_client.HighReadableFunctionRunProgress{ // report progress_milestone
		ProgressMilestoneIndex: ParsingParam.MilestoneIndex(),
	}
	sleepSeconds, err := ipts.GetIntValue(0, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get sleep_seconds from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       false,
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
		}
		return
	}

	logger.Infof("start sleep %d seconds", sleepSeconds)
	progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
		ProgressMilestoneIndex: Sleeping.MilestoneIndex(),
	}
	if sleepSeconds > 0 {
		progressReportTicker := time.NewTicker(time.Second)
		defer progressReportTicker.Stop()
		finishedTicker := time.NewTicker(time.Duration(sleepSeconds) * time.Second)
		defer finishedTicker.Stop()
		sleeppedSeconds := 0
	FOR:
		for {
			select {
			case <-progressReportTicker.C:
				sleeppedSeconds++
				// report log & progress percent every one second
				logger.Infof("sleeped %d/%d seconds", sleeppedSeconds, sleepSeconds)
				progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
					Progress: float32(sleeppedSeconds*100) / float32(sleepSeconds),
				}
			case <-finishedTicker.C:
				break FOR
			}
		}
		progressReportChan <- bloc_client.HighReadableFunctionRunProgress{
			ProgressMilestoneIndex: Finish.MilestoneIndex(),
		}
	}

	logger.Infof("finished")
	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc:                       true,
		InterceptBelowFunctionRun: false,
		Description:               fmt.Sprintf("sleeppedSeconds %d seconds", sleepSeconds),
	}
}
