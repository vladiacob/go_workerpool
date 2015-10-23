package go_workerpool

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	maxWorkers := 2
	workerPool := make(chan chan Job, maxWorkers)
	worker := NewWorker(1, workerPool, true)
	expectedWorker := &Worker{
		id:         1,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		started:    false,
	}

	assert.Equal(t, expectedWorker.id, worker.id)
	assert.Equal(t, expectedWorker.workerPool, worker.workerPool)
	assert.Equal(t, expectedWorker.started, worker.started)
}

func TestStart(t *testing.T) {
	// Create worker
	maxWorkers := 2
	workerPool := make(chan chan Job, maxWorkers)
	worker := NewWorker(1, workerPool, true)
	worker.Start()

	// Create job and add to worker pool
	noJobs := 3
	job := &JobTest{noCalls: 0, response: nil}
	for i := 0; i < noJobs; i++ {
		workerJobQueue := <-workerPool
		workerJobQueue <- job
	}

	// Retry until worker jobs are done
	noRetries := 0
	maxRetries := 10
	testTicker := time.NewTicker(time.Millisecond * 10)
	for _ = range testTicker.C {
		noRetries++
		if job.noCalls == noJobs {
			break
		} else if noRetries == maxRetries {
			break
		}
	}
	testTicker.Stop()
	worker.Stop()

	assert.Equal(t, noJobs, job.noCalls)
}

func TestStartError(t *testing.T) {
	// Create worker
	maxWorkers := 2
	workerPool := make(chan chan Job, maxWorkers)
	worker := NewWorker(1, workerPool, true)
	errMsgs := []error{}
	worker.verbose = func(debug bool, msg string, args ...interface{}) {
		errMsgs = append(errMsgs, fmt.Errorf(msg, args...))
	}
	worker.Start()

	// Create job and add to worker pool
	noJobsSuccess := 3
	noJobsError := 2
	jobSuccess := &JobTest{noCalls: 0, response: nil}
	for i := 0; i < noJobsSuccess; i++ {
		workerJobQueue := <-workerPool
		workerJobQueue <- jobSuccess
	}

	jobError := &JobTest{noCalls: 0, response: errors.New("worker error")}
	for i := 0; i < noJobsError; i++ {
		workerJobQueue := <-workerPool
		workerJobQueue <- jobError
	}

	// Retry until worker jobs are done
	noRetries := 0
	maxRetries := 10
	testTicker := time.NewTicker(time.Millisecond * 10)
	for _ = range testTicker.C {
		noRetries++

		if jobSuccess.noCalls == noJobsSuccess && jobError.noCalls == noJobsError {
			break
		} else if noRetries == maxRetries {
			break
		}
	}
	testTicker.Stop()
	worker.Stop()

	assert.Equal(t, noJobsSuccess, jobSuccess.noCalls)

	assert.Equal(t, noJobsError, jobError.noCalls)
	assert.Equal(t, noJobsError, len(errMsgs))
}

func TestStopAndStarted(t *testing.T) {
	// Create worker
	maxWorkers := 2
	workerPool := make(chan chan Job, maxWorkers)
	worker := NewWorker(1, workerPool, true)
	var infoMsg string
	worker.verbose = func(debug bool, msg string, args ...interface{}) {
		infoMsg = fmt.Sprintf(msg, args...)
	}
	worker.Start()

	job := &JobTest{noCalls: 0, response: nil}
	workerJobQueue := <-workerPool
	workerJobQueue <- job

	worker.Stop()

	// Retry until worker was stoped
	noRetries := 0
	maxRetries := 10
	testTicker := time.NewTicker(time.Millisecond * 10)
	for _ = range testTicker.C {
		noRetries++
		if worker.Started() == false {
			break
		} else if noRetries == maxRetries {
			break
		}
	}

	assert.Equal(t, worker.Started(), false)
	assert.Equal(t, "worker 1 stopping\n", infoMsg)
}

func TestID(t *testing.T) {
	// Create worker
	maxWorkers := 2
	workerPool := make(chan chan Job, maxWorkers)
	worker := NewWorker(1, workerPool, true)
	worker.Start()

	assert.Equal(t, 1, worker.ID())

	worker.Stop()
}

// JobTest structure which simulate a job
type JobTest struct {
	noCalls  int
	workerID int

	response error
}

func (j *JobTest) Work() error {
	j.noCalls++

	return j.response
}

func (j *JobTest) SetWorkerID(ID int) {
	j.workerID = ID
}
