package main

import (
	"log"
	"sync"
)

type Task struct {
	Payload []byte
}

type workerPool struct {
	tasks chan Task
	wg    sync.WaitGroup

	// mu          sync.Mutex for later use // dynamic resizing
	workerCount int
}

func newWorkerPool(numWorkers int) *workerPool {
	pool := &workerPool{
		tasks:       make(chan Task),
		workerCount: numWorkers,
	}

	for range pool.workerCount {
		pool.wg.Go(func() {
			pool.worker()
		})

	}
	return pool
}

func (p *workerPool) worker() {
	for task := range p.tasks {
		// Process the task. In your case, this is where you'd handle the webhook event.
		log.Printf("Processing task with payload: %s\n", string(task.Payload))
	}
}

func (p *workerPool) Submit(task Task) {
	p.tasks <- task
}

func (p *workerPool) Stop() {
	close(p.tasks)
	p.wg.Wait()
}
