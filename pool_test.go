package go_workerpool

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	maxWorkers := 2
	maxJobQueue := 2
	pool := New(maxWorkers, maxJobQueue)

	assert.Equal(t, maxJobQueue, pool.maxJobQueue)
	assert.Equal(t, maxWorkers, pool.maxWorkers)
	assert.Equal(t, NotStarted, pool.status)
}

func TestRun(t *testing.T) {
	maxWorkers := 2
	maxJobQueue := 2
	pool := New(maxWorkers, maxJobQueue)
	pool.Run()

	assert.Equal(t, maxWorkers, len(pool.workers))
	assert.Equal(t, Started, pool.status)
}

func TestAdd(t *testing.T) {
	maxWorkers := 2
	maxJobQueue := 3
	pool := New(maxWorkers, maxJobQueue)
	pool.Run()

	// Create job and add to worker pool
	noJobs := 3
	job := &JobTest{noCalls: 0, response: nil}
	for i := 0; i < noJobs; i++ {
		pool.Add(job)
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

	assert.Equal(t, noJobs, job.noCalls)
}

func TestAddErrors(t *testing.T) {
	maxWorkers := 2
	maxJobQueue := 2
	pool := New(maxWorkers, maxJobQueue)

	job := &JobTest{noCalls: 0, response: nil}
	err := pool.Add(job)

	// Worker pool is not started
	expectedErr := errors.New("job queue have status: not started")

	assert.Equal(t, expectedErr, err)

	pool.Run()

	// Worker pool is full
	err1 := pool.Add(job)
	err2 := pool.Add(job)
	err3 := pool.Add(job)

	expectedErr = errors.New("job queue is full, it have 2 jobs")

	assert.Equal(t, nil, err1)
	assert.Equal(t, nil, err2)
	assert.Equal(t, expectedErr, err3)
}

func TestStatus(t *testing.T) {
	maxWorkers := 1
	maxJobQueue := 1
	pool := New(maxWorkers, maxJobQueue)

	pool.status = NotStarted
	assert.Equal(t, pool.Status(), "not started")

	pool.status = Started
	assert.Equal(t, pool.Status(), "started")

	pool.status = Stopped
	assert.Equal(t, pool.Status(), "stopped")

	pool.status = 4
	assert.Equal(t, pool.Status(), "unknown")
}

func TestWaitAndStop(t *testing.T) {
	maxWorkers := 2
	maxJobQueue := 10
	pool := New(maxWorkers, maxJobQueue)
	pool.Run()

	// Create job and add to worker pool
	noJobs := 10
	job := &JobWaitTest{noCalls: 0, response: nil}
	for i := 0; i < noJobs; i++ {
		pool.Add(job)
	}

	err := pool.Stop(true)

	assert.Equal(t, noJobs, job.noCalls)
	assert.Nil(t, err)
}

func TestImmediateStop(t *testing.T) {
	maxWorkers := 1
	maxJobQueue := 5
	pool := New(maxWorkers, maxJobQueue)
	pool.Run()

	// Create job and add to worker pool
	noJobs := 5
	job := &JobWaitTest{noCalls: 0, response: nil}
	for i := 0; i < noJobs; i++ {
		pool.Add(job)
	}

	pool.Stop(false)

	assert.NotEqual(t, noJobs, job.noCalls)
}

// JobTest structure which simulate a job
type JobWaitTest struct {
	noCalls int

	response error
}

func (j *JobWaitTest) Work() error {
	time.Sleep(1 * time.Second)

	j.noCalls++

	return j.response
}
