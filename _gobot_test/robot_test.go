package gobot_test

import (
	"errors"
	"testing"

	"sync"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot"
)

// This test shows that a Robot set up with a panicHandler function
// can set its lastError field and be visible to the caller
func TestRobotPanicHandler_SetsLastError(t *testing.T) {
	signalErrorValue := errors.New("captured error")

	panicHandler := func(r *gobot.Robot, err error) {
		r.SetLastError(err)
		r.Stop()
	}

	panicConfigFunc := func(r *gobot.Robot) {
		r.SetPanicHandler(panicHandler)
	}

	panickyWork := func() {
		panic(signalErrorValue)
	}

	r := gobot.NewRobot(panickyWork)
	r.AutoRun = false

	r.AddConfigFunc(panicConfigFunc)
	r.RunConfigFuncs()
	r.SetWaitGroup(&sync.WaitGroup{})

	r.Start()
	r.WaitGroup().Wait()

	assert.Equal(t, signalErrorValue, r.LastError())
}
