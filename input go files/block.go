package main

import (
	"math"
	"time"
)

func goBlocked(ch chan int) {
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10000000; i++ {
		j := math.Sqrt(float64(i * i))
		j *= j * float64(i)
		_ = j
	}
	val := 42
	ch <- val

	time.Sleep(10 * time.Millisecond)
}

func main() {
	time.Sleep(10 * time.Millisecond)
	ch := make(chan int)
	go goBlocked(ch)
	time.Sleep(100 * time.Millisecond)
}
