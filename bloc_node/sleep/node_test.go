package sleep

import (
	"testing"
	"time"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func TestSend(t *testing.T) {
	client := bloc_client.NewTestClient()
	toSleepSeconds := 10
	beforeSleep := time.Now()
	executeOpt := client.TestRunFunction(
		&Sleep{},
		[][]interface{}{
			{
				toSleepSeconds,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should fail, but suc")
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept_below_functionrun")
	}
	if !time.Now().Add(-time.Duration(toSleepSeconds) * time.Second).After(beforeSleep) {
		t.Fatal("sleep time not enough")
	}
}
