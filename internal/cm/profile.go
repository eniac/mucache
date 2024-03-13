package cm

import (
	"fmt"
	"time"
)

type ProfileState struct {
	Counter       int
	Time          time.Time
	PrintTimeFreq int
}

func initProfile(printTimeFreq int) *ProfileState {
	ret := ProfileState{
		Counter:       0,
		Time:          time.Now(),
		PrintTimeFreq: printTimeFreq,
	}
	return &ret
}

func (state *ProfileState) profileProcRequest() {
	state.Counter += 1
	if state.Counter == state.PrintTimeFreq {
		state.Counter = 0
		t := time.Now()
		timeDiff := t.Sub(state.Time)
		state.Time = t
		fmt.Printf("Processed %v reqeusts in %v\n", state.PrintTimeFreq, timeDiff.Microseconds())
	}
}
