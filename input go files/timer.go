package main

import (
	"time"
)

func goTimer(d time.Duration, c chan int) {
	time.Sleep(d)
	c <- 1
	time.Sleep(time.Millisecond * 300)
}

func timer(d time.Duration) <-chan int {
	c := make(chan int)
	go goTimer(d, c)
	return c
}

func main() {
	for i := 0; i < 24; i++ {
		c := timer(1 * time.Second)
		<-c
	}
}
