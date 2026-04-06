package handler

import "sync"

type WorkerPool struct {
	tasks chan func()
	wg    sync.WaitGroup
	size  int
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		tasks: make(chan func(), 100),
		size:  size,
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.tasks {
				task()
			}
		}()
	}
}

func (p *WorkerPool) Submit(task func()) {
	p.tasks <- task
}

func (p *WorkerPool) Stop() {
	close(p.tasks)
	p.wg.Wait()
}
