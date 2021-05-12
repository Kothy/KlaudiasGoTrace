package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"fmt"
	"sync"
	"time"
)

func worker2(tasksCh <-chan int, wg *sync.WaitGroup, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	defer wg.Done()
	for {
		task, ok := <-tasksCh
		KlaudiasGoTrace.ReceiveFromChannel(task,
			tasksCh)

		if !ok {
			return
		}
		d := time.Duration(task) * time.Millisecond
		time.Sleep(d)
		fmt.Println("processing task", task)
	}
	KlaudiasGoTrace.StopGoroutine()

}

func pool(wg *sync.WaitGroup, workers, tasks int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	tasksCh := make(chan int)

	for i := 0; i < workers; i++ {
		go worker2(tasksCh, wg, KlaudiasGoTrace.GetGID())
	}

	for i := 0; i < tasks; i++ {
		KlaudiasGoTrace.SendToChannel(
			i, tasksCh)

		tasksCh <- i
	}

	close(tasksCh)
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	var wg sync.WaitGroup
	wg.Add(36)
	go pool(&wg, 36, 50, KlaudiasGoTrace.GetGID())
	wg.Wait()
	KlaudiasGoTrace.EndTrace()

}
