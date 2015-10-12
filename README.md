# go_workerpool
Generic worker pool for Go Language. You can limit the number of go routines which are running in parallel.

## Godoc
[https://godoc.org/github.com/vladiacob/go_workerpool](https://godoc.org/github.com/vladiacob/go_workerpool)

## How to install
```
go get github.com/vladiacob/go_workerpool
```

## How to run examples
```
cd examples
go run example.go
```

## How to use
### Include go_workerpool
```
include (
    ..
    pool "github.com/vladiacob/go_workerpool"
    ..
)
```

### Initialize pool
```
workerPool := pool.New(2, 10)
workerPool.Run()
```
* maxWorkers: number of workers which are processing in parallel
* maxJobQueue: number of jobs which will be accept in queue

### Add job to pool
```
err := workerPool.Add(job1)
if err != nil {
    fmt.Println(err)
}
```

### Stop worker pool
```
err := workerPool.Stop(true)
if err != nil {
    fmt.Println(err)
}
```
* waitAndStop == true: wait until all jobs was processed
* waitAndStop == false: close workers immediate 
