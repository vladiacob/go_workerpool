package go_workerpool

import (
	"errors"
	"fmt"
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

	// Worker pool is full
	maxWorkers = 1
	maxJobQueue = 1
	pool = New(maxWorkers, maxJobQueue)
	pool.Run()

	jobWait := &JobWaitTest{noCalls: 0, response: nil}
	err1 := pool.Add(jobWait)
	time.Sleep(10 * time.Microsecond)
	err2 := pool.Add(jobWait)
	time.Sleep(10 * time.Microsecond)
	err3 := pool.Add(jobWait)
	time.Sleep(10 * time.Microsecond)
	err4 := pool.Add(jobWait)

	expectedErr = errors.New("job queue is full, it have 1 jobs")
	assert.Equal(t, nil, err1)
	assert.Equal(t, nil, err2)
	assert.Equal(t, nil, err3)
	assert.Equal(t, expectedErr, err4)
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

func TestStats(t *testing.T) {
	maxWorkers := 5
	maxJobQueue := 10
	pool := New(maxWorkers, maxJobQueue)
	pool.Run()

	// Create job and add to worker pool
	noJobs := 4

	for i := 0; i < noJobs; i++ {
		job := &JobWaitTest{noCalls: 0, response: nil, name: i}
		pool.Add(job)
	}

	stats := pool.Stats()
	assert.Equal(t, 6, stats["free_job_queue_spaces"])

	time.Sleep(1 * time.Second)

	stats = pool.Stats()
	assert.Equal(t, 1, stats["free_workers"])

	pool.Stop(true)
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
	noCalls  int
	workerID int
	name     int

	response error
}

func (j *JobWaitTest) Work() error {
	time.Sleep(5 * time.Second)
	fmt.Println(fmt.Sprintf("test %d", j.workerID))
	j.noCalls++

	return j.response
}

func (j *JobWaitTest) SetWorkerID(ID int) {
	j.workerID = ID
}

func (j *JobWaitTest) Name() int {
	return j.name
}
