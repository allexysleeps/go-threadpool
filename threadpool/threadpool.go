package threadpool

import (
	"context"
	"sync"
	"sync/atomic"
)

func dummyID() func() uint8 {
	mu := sync.Mutex{}
	id := uint8(0)
	return func() uint8 {
		mu.Lock()
		newID := id
		id++
		mu.Unlock()
		return newID
	}
}

const (
	StatusPending  = "Pending"
	StatusRunning  = "Running"
	StatusDone     = "Done"
	StatusCanceled = "Canceled"
)

type Operation func(ctx context.Context)

type Threadpool struct {
	size          int
	activeWorkers int32
	queue         []*Task
	workload      chan *Task
	ctx           context.Context
	genID         func() uint8
	sync.RWMutex
}

type Task struct {
	tp        *Threadpool
	id        uint8
	status    string
	operation Operation
	ctx       context.Context
	cancel    context.CancelFunc
}

func (tp *Threadpool) Run(op Operation) *Task {
	ctx, cancel := context.WithCancel(tp.ctx)
	task := Task{
		id:        tp.genID(),
		operation: op,
		status:    StatusPending,
		ctx:       ctx,
		cancel:    cancel,
		tp:        tp,
	}

	tp.Lock()
	tp.queue = append(tp.queue, &task)
	tp.Unlock()

	return &task
}

func (t *Task) Status() string {
	t.tp.RLock()
	status := t.status
	t.tp.RUnlock()
	return status
}

func (t *Task) Stop() {
	t.tp.Lock()
	defer t.tp.Unlock()

	switch t.status {
	case StatusRunning:
		t.cancel()
		t.status = StatusCanceled
		return
	case StatusPending:
		for i, tsk := range t.tp.queue {
			if tsk.id != t.id {
				continue
			}
			if i == len(t.tp.queue)-1 {
				t.tp.queue = t.tp.queue[:i]
			}
			t.tp.queue = append(t.tp.queue[0:i], t.tp.queue[i+1:len(t.tp.queue)]...)

			t.status = StatusCanceled
			return
		}
	}
}

func startWorker(tp *Threadpool) {
	for {
		select {
		case task := <-tp.workload:
			task.status = StatusRunning
			task.operation(task.ctx)
			if task.status != StatusCanceled {
				task.status = StatusDone
			}
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
		ctx:      context.Background(),
		genID:    dummyID(),
	}

	for i := 0; i < count; i++ {
		go startWorker(&tp)
	}

	return &tp
}
