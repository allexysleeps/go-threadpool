package threadpool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	plannedOp := 100
	performedOp := 0

	tp := Create(4)

	for i := 0; i < plannedOp; i++ {
		func() {
			tp.Run(func(ctx context.Context) {
				time.Sleep(0)
				performedOp++
			})
		}()
		fmt.Println(i)
	}
loop:
	for {
		select {
		case <-time.After(time.Second):
			t.Errorf("threadpool hasnt finished it 5s")
			break loop
		default:
			if plannedOp == performedOp {
				break loop
			}
		}
	}
}
