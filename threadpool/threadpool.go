package threadpool

import (
	"sync"
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
	activeWorkers int
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

	if tp.activeWorkers == tp.size {
		task.status = StatusPending
		tp.Lock()
		tp.queue = append(tp.queue, &task)
		tp.Unlock()
		return &task
	}

	task.status = StatusRunning

	tp.Lock()
	tp.activeWorkers++
	tp.Unlock()

	tp.workload <- &task

	return &task
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func Create(count int) *Threadpool {
	tp := Threadpool{
		size:     count,
		queue:    make([]*Task, 0),
		workload: make(chan *Task),
	}

	for i := 0; i < count; i++ {
		go func() {
			for {
				select {
				case task := <-tp.workload:
					task.status = StatusRunning
					task.operation()
					task.status = StatusDone

					tp.Lock()
					tp.activeWorkers--
					tp.Unlock()
				default:
					var task *Task
					tp.Lock()

					if len(tp.queue) > 0 {
						task = tp.queue[0]
						tp.activeWorkers++
						if len(tp.queue) == 1 {
							tp.queue = []*Task{}
						} else {
							tp.queue = tp.queue[1:len(tp.queue)]
						}
					}
					tp.Unlock()
					if task != nil {
						tp.workload <- task
					}
				}
			}
		}()
	}

	return &tp
}
