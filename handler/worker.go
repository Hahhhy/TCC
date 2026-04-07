package handler

import "sync"


//池，WorkerPool是一个简单的协程池实现，
// 维护一个任务队列和固定数量的Worker协程来处理任务。、
// Submit方法用于提交任务，Stop方法用于关闭池并等待所有Worker完成。
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
