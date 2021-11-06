# threadpool-go

Simple thread pool implementation in go

```go
tp := threadpool.Create(N) // N - amount of goroutines needs to be running
task := tp.Run(func() {
	time.Sleep(time.Seconds * 5)
})
task.Status() // Running
```