package phone_sms

import (
	"testing"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func TestSend(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&SMS{},
		[][]interface{}{
			{
				"xx",
			},
			{
				[]string{"11111111", "00000000"},
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should fail, but suc")
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept_below_functionrun")
	}
}
