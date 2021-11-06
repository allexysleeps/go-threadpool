package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/threadpool/threadpool"
)

func main() {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	plannedOp := 1000
	performedOp := 0

	var mu sync.Mutex

	tp := threadpool.Create(4)

	for i := 0; i < plannedOp; i++ {
		func(idx int) {
			tp.Run(func() {
				delay := r.Intn(3)
				time.Sleep(time.Second * time.Duration(delay))

				mu.Lock()
				performedOp = performedOp + 1
				mu.Unlock()
				log.Printf("operation #%d finished, delay %d\n", idx, delay)
			})
		}(i)
	}

	fmt.Println("--------loaded chunk---------")

	for performedOp != plannedOp {
	}
}
