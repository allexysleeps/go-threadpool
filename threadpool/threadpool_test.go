package threadpool

import (
	"sync"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	plannedOp := 10000
	performedOp := 0
	var mu sync.Mutex

	tp := Create(5)

	for i := 0; i < plannedOp; i++ {
		func() {
			tp.Run(func() {
				mu.Lock()
				performedOp++
				mu.Unlock()
			})
		}()
	}
loop:
	for {
		select {
		case <-time.After(time.Second * 5):
			t.Errorf("treahdpool hasnt finished it 5s")
		default:
			if plannedOp == performedOp {
				break loop
			}
		}
	}
}
