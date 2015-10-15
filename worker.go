package go_workerpool

import (
	"log"
)

// Job interface which will be used to create a new job
type Job interface {
	Work() error
}

// Worker is the structure for worker
type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
	started    bool
	debug      bool
	verbose    func(debug bool, msg string, args ...interface{})
}

// NewWorker return a new instance of worker
func NewWorker(id int, workerPool chan chan Job, debug bool) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		started:    false,
		debug:      debug,
		verbose:    verbose,
	}
}

// Start worker
func (w *Worker) Start() {
	w.started = true

	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				if err := job.Work(); err != nil {
					w.verbose(w.debug, "error running worker %d: %s\n", w.id, err.Error())
				}

			case <-w.quitChan:
				w.verbose(w.debug, "worker %d stopping\n", w.id)

				w.started = false

				return
			}
		}
	}()
}

// Stop worker
func (w *Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}

// ID return worker id
func (w *Worker) ID() int {
	return w.id
}

// Started return worker status
func (w *Worker) Started() bool {
	return w.started
}

func verbose(debug bool, msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}
