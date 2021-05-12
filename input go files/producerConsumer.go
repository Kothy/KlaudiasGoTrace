package main

import (
	"time"
)

func producer(ch chan int, d time.Duration) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(d)
	}
}

func consumer(out chan int) {
	for i := 0; i < 20; i++ {
		<-out
	}
}

func main() {
	ch := make(chan int)
	out := make(chan int)

	go producer(ch, 10*time.Millisecond)
	go producer(ch, 25*time.Millisecond)
	go consumer(out)

	for i := 0; i < 20; i++ {
		i := <-ch
		out <- i
	}
}
