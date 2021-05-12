package main

import (
	"time"
)

func goHello(ch chan int) {
	time.Sleep(10 * time.Millisecond)
	val := 42
	ch <- val
	time.Sleep(10 * time.Millisecond)
}

func main() {
	ch := make(chan int)

	go goHello(ch)

	<-ch
}
