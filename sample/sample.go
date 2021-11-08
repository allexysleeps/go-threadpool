package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/threadpool/threadpool"
)

func basicUsage() {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	plannedOp := 1000
	performedOp := 0

	var mu sync.Mutex

	tp := threadpool.Create(4)

	for i := 0; i < plannedOp; i++ {
		func(idx int) {
			tp.Run(func(ctx context.Context) {
				delay := r.Intn(3)
				time.Sleep(time.Second * time.Duration(delay))
				ctx.Done()

				mu.Lock()
				performedOp = performedOp + 1
				mu.Unlock()
				log.Printf("operation #%d finished, delay %d\n", idx, delay)
			})
		}(i)
	}

	log.Println("--------loaded chunk---------")

	for performedOp != plannedOp {
	}
}

func canceledOperation() {
	done := make(chan bool)
	tp := threadpool.Create(2)
	tsk := tp.Run(func(ctx context.Context) {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				log.Printf("operation canceled")
				close(done)
				return
			case <-ticker.C:
				log.Print("operation is running")
			}
		}
	})

	log.Printf("task status %s\n", tsk.Status())

	go func() {
		time.Sleep(time.Second * 2)
		log.Printf("task status %s\n", tsk.Status())
		tsk.Stop()
	}()

	<-done
	log.Printf("task status %s\n", tsk.Status())
}

func main() {
	canceledOperation()
}
