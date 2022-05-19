package sms

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var failReasons = []string{
	"phone number not exist",                 // 号码不存在
	"phone number already run out of credit", // 欠费
	"remote phone reject",                    //对方拒绝接收
}

func pickARandomFailError() error {
	return errors.New(failReasons[rand.Intn(len(failReasons))])
}

var (
	sucOrFail bool
	lock      sync.Mutex
)

func SendMsg(phoneNumber, msg string) error {
	// we won't really implement this function.
	// Alternating suc or fail(with fail reason) send
	lock.Lock()
	defer lock.Unlock()
	sendSuc := !sucOrFail
	sucOrFail = sendSuc

	if sendSuc {
		return nil
	}
	return pickARandomFailError()
}
