package go_workerpool

import (
	"fmt"
	"time"
)

const (
	// NotStarted is status for not started
	NotStarted int = 0
	// Started is status for started
	Started int = 1
	// Stopped is status for stopped
	Stopped int = 2
)

// Pool is the structure for worker pool
type Pool struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool  chan chan Job
	jobQueue    chan Job
	maxWorkers  int
	maxJobQueue int
	workers     []*Worker
	status      int
}

// New return a new instance of pool
func New(maxWorkers int, maxJobQueue int) *Pool {
	// Create the job queue
	jobQueue := make(chan Job, maxJobQueue)

	// Create the worker pool
	workerPool := make(chan chan Job, maxWorkers)

	return &Pool{
		workerPool:  workerPool,
		jobQueue:    jobQueue,
		maxJobQueue: maxJobQueue,
		maxWorkers:  maxWorkers,
		status:      NotStarted,
	}
}

// Run start all workers
func (p *Pool) Run() {
	// starting n number of workers
	for i := 0; i < p.maxWorkers; i++ {
		worker := NewWorker(i+1, p.workerPool, true)
		worker.Start()
		p.workers = append(p.workers, worker)
	}

	go p.dispatch()

	p.status = Started
}

// Add a new job to be processed
func (p *Pool) Add(job Job) error {
	if p.status == Stopped || p.status == NotStarted {
		return fmt.Errorf("job queue have status: %s", p.Status())
	}

	jobQueueSize := len(p.jobQueue)
	if jobQueueSize == p.maxJobQueue {
		return fmt.Errorf("job queue is full, it have %d jobs", jobQueueSize)
	}

	p.jobQueue <- job

	return nil
}

// Stop worker pool
func (p *Pool) Stop(waitAndStop bool) error {
	p.status = Stopped

	if waitAndStop {
		return p.waitAndStop()
	}

	return p.immediateStop()
}

// Status return worker pool status
func (p *Pool) Status() string {
	switch {
	case p.status == NotStarted:
		return "not started"
	case p.status == Started:
		return "started"
	case p.status == Stopped:
		return "stopped"
	}

	return "unknown"
}

func (p *Pool) waitAndStop() error {
	retryInterval := time.Duration(500)

	retriesLimit := int(time.Duration(60000) / retryInterval)
	retries := 0

	queueTicker := time.NewTicker(time.Millisecond * retryInterval)
	// Retry evenry 500ms to check if job queue is empty
	for _ = range queueTicker.C {
		// Check if jobQueue is empty and all workers are available
		if len(p.jobQueue) == 0 && len(p.workerPool) == p.maxWorkers {
			for _, worker := range p.workers {
				worker.Stop()
			}

			break
		}

		retries++
		if retries >= retriesLimit {
			queueTicker.Stop()

			return fmt.Errorf(fmt.Sprintf("checking queue status exceeded retry limit: %v", time.Duration(retries)*retryInterval*time.Millisecond))
		}
	}

	// Stop job queue ticker
	queueTicker.Stop()

	return nil
}

func (p *Pool) immediateStop() error {
	for _, worker := range p.workers {
		worker.Stop()
	}

	return nil
}

func (p *Pool) dispatch() {
	for {
		select {
		case job := <-p.jobQueue:
			go func() {
				workerJobQueue := <-p.workerPool
				workerJobQueue <- job
			}()
		}
	}
}
