package main

import (
	"fmt"
	"time"

	pool "github.com/vladiacob/go_workerpool"
)

// Job structure
type Job struct {
	name     string
	workerID int
}

// Work is method which is called by the worker
func (j *Job) Work() error {
	fmt.Println("Work: " + j.name)

	time.Sleep(1 * time.Second)

	return nil
}

func (j *Job) SetWorkerID(ID int) {
	j.workerID = ID
}

func main() {
	// Initializate worker pool
	workerPool := pool.New(2, 10)
	workerPool.Run()

	// Create 10 jobs
	job1 := &Job{name: "job1"}
	job2 := &Job{name: "job2"}
	job3 := &Job{name: "job3"}
	job4 := &Job{name: "job4"}
	job5 := &Job{name: "job5"}
	job6 := &Job{name: "job6"}
	job7 := &Job{name: "job7"}
	job8 := &Job{name: "job8"}
	job9 := &Job{name: "job9"}
	job10 := &Job{name: "job10"}

	// Add jobs to worker pool
	go func() {
		err := workerPool.Add(job1)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job2)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job3)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job4)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job5)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job6)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job7)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job8)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job9)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := workerPool.Add(job10)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Wait
	time.Sleep(2 * time.Second)

	// Stop worker pool and wait until all jobs are done
	fmt.Println("Stopping...")
	err := workerPool.Stop(true)
	if err != nil {
		fmt.Println(err)
	}

	// Try to add a new job (will return error)
	go func() {
		err := workerPool.Add(job10)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Wait
	time.Sleep(20 * time.Second)
}
