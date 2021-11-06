package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/threadpool/threadpool"
)

func main() {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	plannedOp := 50
	performedOp := 0

	var mu sync.Mutex

	tp := threadpool.Create(5)

	for i := 0; i < plannedOp; i++ {
		func(idx int) {
			tp.Run(func() {
				delay := r.Intn(5)
				time.Sleep(time.Second * time.Duration(delay))
				fmt.Printf("operation #%d done\n", idx)

				mu.Lock()
				performedOp = performedOp + 1
				mu.Unlock()
			})
		}(i)
	}

	fmt.Println("--------loaded chunk---------")

	for performedOp != plannedOp {
	}
}
