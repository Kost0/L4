package someWork

import (
	"math/rand"
	"runtime/debug"
	"time"
)

func StartWork() {
	timer := time.NewTimer(30 * time.Second)

	for {
		<-timer.C

		var sl []int

		for i := range 1000 {
			sl = append(sl, i*5/2*3)
		}

		debug.SetGCPercent(rand.Intn(100))

		timer.Reset(30 * time.Second)
	}
}
