package threadpool

import (
	"sync"
	"sync/atomic"
)

type TaskStatus uint8

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusDone
)

type operation func()

type Threadpool struct {
	size          int
	activeWorkers int32
	queue         []*Task
	workload      chan *Task
	sync.RWMutex
}

type Task struct {
	status    TaskStatus
	operation operation
}

func (tp *Threadpool) Run(op operation) *Task {
	task := Task{
		operation: op,
	}

	tp.Lock()
	defer tp.Unlock()
	if tp.activeWorkers == int32(tp.size) {
		task.status = StatusPending
		tp.queue = append(tp.queue, &task)
		return &task
	}

	tp.activeWorkers++

	task.status = StatusRunning
	tp.workload <- &task

	return &task
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func startWorker(tp *Threadpool) {
	for {
		select {
		case task := <-tp.workload:
			task.status = StatusRunning
			task.operation()
			task.status = StatusDone

			atomic.AddInt32(&tp.activeWorkers, -1)
		default:
			tp.Lock()
			if len(tp.queue) > 0 {
				task := tp.queue[0]
				tp.activeWorkers++
				if len(tp.queue) == 1 {
					tp.queue = tp.queue[:0]
				} else {
					tp.queue = tp.queue[1:len(tp.queue)]
				}
				tp.workload <- task
			}
			tp.Unlock()
		}
	}
}

func Create(count int) *Threadpool {
	tp := Threadpool{
		size:     count,
		queue:    make([]*Task, 0),
		workload: make(chan *Task, count),
	}

	for i := 0; i < count; i++ {
		go startWorker(&tp)
	}

	return &tp
}
