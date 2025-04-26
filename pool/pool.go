package pool

import (
	"context"
)

type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc
	tasks  chan func()
}

func New(size int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		ctx:    ctx,
		tasks:  make(chan func(), size),
		cancel: cancel,
	}

	for i := 0; i < size; i++ {
		go p.work()
	}

	return p
}

func (p *Pool) Stop() {
	p.cancel()
	close(p.tasks)
}

func (p *Pool) Put(task func()) {
	p.tasks <- task
}

func (p *Pool) work() {
	for {
		select {
		case task := <-p.tasks:
			task()
		case <-p.ctx.Done():
			return
		}
	}
}
