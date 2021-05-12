package main

import (
	"time"
)

func goHelloBlocked(ch chan int) {
	time.Sleep(10 * time.Millisecond)
	val := 42
	ch <- val
}

func main() {
	ch := make(chan int)
	go goHelloBlocked(ch)
	<-ch
}
