package threadpool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	plannedOp := 100
	var performedOp int32

	tp := Create(4)

	for i := 0; i < plannedOp; i++ {
		func() {
			tp.Run(func(ctx context.Context) {
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&performedOp, 1)
			})
		}()
	}
loop:
	for {
		select {
		case <-time.After(time.Second):
			t.Errorf("threadpool hasnt finished it 5s")
			break loop
		default:
			if plannedOp == int(performedOp) {
				break loop
			}
		}
	}
}

func TestThreadpool_Run(t *testing.T) {
	t.Parallel()

	tp := Create(2)
	done := make(chan struct{})

	wg := sync.WaitGroup{}

	go func() {
		select {
		case <-time.NewTicker(time.Second * 6).C:
			t.Errorf("Operation took too long")
			done <- struct{}{}
		case <-done:
		}
	}()

	wg.Add(2)
	tp.Run(func(ctx context.Context) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 1)
		}
	})

	tp.Run(func(ctx context.Context) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 1)
		}
	})

	wg.Wait()

	done <- struct{}{}
}
