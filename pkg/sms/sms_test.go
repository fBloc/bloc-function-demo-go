package sms

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(SendMsg("18888888888", "test"))
	}
}
